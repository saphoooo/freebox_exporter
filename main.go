package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mafreebox string
	listen    string
	debug     bool
	fiber     bool
)

func init() {
	flag.StringVar(&mafreebox, "endpoint", "http://mafreebox.freebox.fr/", "Endpoint for freebox API")
	flag.StringVar(&listen, "listen", ":10001", "Prometheus metrics port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&fiber, "fiber", false, "Turn on if you're using a fiber Freebox")
}

func main() {
	flag.Parse()

	if !strings.HasSuffix(mafreebox, "/") {
		mafreebox = mafreebox + "/"
	}

	// myAuthInfo contains all auth data
	endpoint := mafreebox + "api/v4/login/"
	myAuthInfo := &authInfo{
		myAPI: api{
			login:        endpoint,
			authz:        endpoint + "authorize/",
			loginSession: endpoint + "session/",
		},
		myStore: store{location: os.Getenv("HOME") + "/.freebox_token"},
		myApp: app{
			AppID:      "fr.freebox.exporter",
			AppName:    "prometheus-exporter",
			AppVersion: "0.4",
			DeviceName: "local",
		},
		myReader: bufio.NewReader(os.Stdin),
	}

	myPostRequest := newPostRequest()

	myConnectionXdslRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/connection/xdsl/",
		header: "X-Fbx-App-Auth",
	}

	myFreeplugRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/freeplug/",
		header: "X-Fbx-App-Auth",
	}

	myLanRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/lan/browser/pub/",
		header: "X-Fbx-App-Auth",
	}

	mySystemRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/v4/system/",
		header: "X-Fbx-App-Auth",
	}

	var mySessionToken string

	// infinite loop to get all statistics
	go func() {
		for {
			// There is no DSL metric on fiber Freebox
			// If you use a fiber Freebox, use -fiber flag to turn off this metric
			if !fiber {
				// connectionXdsl metrics
				connectionXdslStats, err := getConnectionXdsl(myAuthInfo, myConnectionXdslRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with connectionXdsl metrics: %v", err)
				}

				if connectionXdslStats.Success {
					status := connectionXdslStats.Result.Status
					result := connectionXdslStats.Result
					down := result.Down
					up := result.Up

					connectionXdslStatusUptimeGauge.Set(float64(status.Uptime))

					connectionXdslDownAttnGauge.Set(float64(down.Attn10) / 10)
					connectionXdslUpAttnGauge.Set(float64(up.Attn10) / 10)

					connectionXdslDownSnrGauge.Set(float64(down.Snr10) / 10)
					connectionXdslUpSnrGauge.Set(float64(up.Snr10) / 10)

					resultReflect := reflect.ValueOf(result)
					for _, direction := range []string{"down", "up"} {
						for _, errorField := range []string{"crc", "es", "fec", "hec", "ses"} {
							value := reflect.Indirect(resultReflect).
								FieldByName(strings.Title(direction)).
								FieldByName(strings.Title(errorField))

							connectionXdslErrorGauges.WithLabelValues(direction, errorField).
								Set(float64(value.Int()))
						}
					}
				}

				// dsl metrics
				getDslResult, err := getDsl(myAuthInfo, myPostRequest, &mySessionToken)
				if err != nil {
					log.Printf("An error occured with DSL metrics: %v", err)
				}

				if len(getDslResult) > 0 {
					rateUpGauge.Set(float64(getDslResult[0]))
					rateDownGauge.Set(float64(getDslResult[1]))
					snrUpGauge.Set(float64(getDslResult[2]))
					snrDownGauge.Set(float64(getDslResult[3]))
				}
			}

			// freeplug metrics
			freeplugStats, err := getFreeplug(myAuthInfo, myFreeplugRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with freeplug metrics: %v", err)
			}

			for _, freeplugNetwork := range freeplugStats.Result {
				for _, freeplugMember := range freeplugNetwork.Members {
					if freeplugMember.HasNetwork {
						freeplugHasNetworkGauge.WithLabelValues(freeplugMember.ID).Set(float64(1))
					} else {
						freeplugHasNetworkGauge.WithLabelValues(freeplugMember.ID).Set(float64(0))
					}

					Mb := 1e6
					rxRate := float64(freeplugMember.RxRate) * Mb
					txRate := float64(freeplugMember.TxRate) * Mb

					if rxRate >= 0 { // -1 if not unavailable
						freeplugRxRateGauge.WithLabelValues(freeplugMember.ID).Set(rxRate)
					}

					if txRate >= 0 { // -1 if not unavailable
						freeplugTxRateGauge.WithLabelValues(freeplugMember.ID).Set(txRate)
					}
				}
			}

			// net metrics
			getNetResult, err := getNet(myAuthInfo, myPostRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with NET metrics: %v", err)
			}

			if len(getNetResult) > 0 {
				bwUpGauge.Set(float64(getNetResult[0]))
				bwDownGauge.Set(float64(getNetResult[1]))
				netRateUpGauge.Set(float64(getNetResult[2]))
				netRateDownGauge.Set(float64(getNetResult[3]))
				vpnRateUpGauge.Set(float64(getNetResult[4]))
				vpnRateDownGauge.Set(float64(getNetResult[5]))
			}

			// lan metrics
			lanAvailable, err := getLan(myAuthInfo, myLanRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with LAN metrics: %v", err)
			}
			for _, v := range lanAvailable {
				if v.Reachable {
					lanReachableGauges.WithLabelValues(v.PrimaryName).Set(float64(1))
				} else {
					lanReachableGauges.WithLabelValues(v.PrimaryName).Set(float64(0))
				}
			}

			// system metrics
			systemStats, err := getSystem(myAuthInfo, mySystemRequest, &mySessionToken)
			if err != nil {
				log.Printf("An error occured with System metrics: %v", err)
			}

			systemTempGauges.WithLabelValues("Température CPU B").Set(float64(systemStats.Result.TempCpub))
			systemTempGauges.WithLabelValues("Température CPU M").Set(float64(systemStats.Result.TempCpum))
			systemTempGauges.WithLabelValues("Disque dur").Set(float64(systemStats.Result.TempHDD))
			systemTempGauges.WithLabelValues("Température Switch").Set(float64(systemStats.Result.TempSW))
			systemFanGauges.WithLabelValues("Ventilateur 1").Set(float64(systemStats.Result.FanRPM))

			time.Sleep(10 * time.Second)
		}
	}()

	// expose the registered metrics via HTTP OpenMetrics
	log.Println("freebox_exporter started on port", listen)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
