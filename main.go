package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mafreebox string
	listen    string
)

func init() {
	flag.StringVar(&mafreebox, "endpoint", "http://mafreebox.freebox.fr/", "Endpoint for freebox API")
	flag.StringVar(&listen, "listen", ":10001", "prometheus metrics port")
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
			// dsl metrics
			getDslResult, err := getDsl(myAuthInfo, myPostRequest, &mySessionToken)
			if err != nil {
				log.Print(err)
			}

			if len(getDslResult) == 0 {
				rateUpGauge.Set(float64(0))
				rateDownGauge.Set(float64(0))
				snrUpGauge.Set(float64(0))
				snrDownGauge.Set(float64(0))
			} else {
				rateUpGauge.Set(float64(getDslResult[0]))
				rateDownGauge.Set(float64(getDslResult[1]))
				snrUpGauge.Set(float64(getDslResult[2]))
				snrDownGauge.Set(float64(getDslResult[3]))
			}

			// switch metrcis
			// as switch database seems to be broken, this one is not used at this time
			/*
				getSwitchResult, err := getSwitch(myAuthInfo, myPostRequest, &mySessionToken)
				if err != nil {
					log.Print(err)
				}

				if len(getSwitchResult) == 0 {
					rx1Gauge.Set(float64(0))
					tx1Gauge.Set(float64(0))
					rx2Gauge.Set(float64(0))
					tx2Gauge.Set(float64(0))
					rx3Gauge.Set(float64(0))
					tx3Gauge.Set(float64(0))
					rx4Gauge.Set(float64(0))
					tx4Gauge.Set(float64(0))
				} else {
					rx1Gauge.Set(float64(getSwitchResult[0]))
					tx1Gauge.Set(float64(getSwitchResult[1]))
					rx2Gauge.Set(float64(getSwitchResult[2]))
					tx2Gauge.Set(float64(getSwitchResult[3]))
					rx3Gauge.Set(float64(getSwitchResult[4]))
					tx3Gauge.Set(float64(getSwitchResult[5]))
					rx4Gauge.Set(float64(getSwitchResult[6]))
					tx4Gauge.Set(float64(getSwitchResult[7]))
				}
			*/

			// temps metrcis
			// as temp database seems to be broken, this one is not used at this time
			// system report the same kind of value
			/*
				getTempResult, err := getTemp(myAuthInfo, myPostRequest, &mySessionToken)
				if err != nil {
					log.Print(err)
				}

				if len(getTempResult) == 0 {
					cpumGauge.Set(float64(0))
					cpubGauge.Set(float64(0))
					swGauge.Set(float64(0))
					hddGauge.Set(float64(0))
					fanSpeedGauge.Set(float64(0))
				} else {
					cpumGauge.Set(float64(getTempResult[0]))
					cpubGauge.Set(float64(getTempResult[1]))
					swGauge.Set(float64(getTempResult[2]))
					hddGauge.Set(float64(getTempResult[3]))
					fanSpeedGauge.Set(float64(getTempResult[4]))
				}
			*/

			// net metrics
			getNetResult, err := getNet(myAuthInfo, myPostRequest, &mySessionToken)
			if err != nil {
				log.Print(err)
			}

			if len(getNetResult) == 0 {
				bwUpGauge.Set(float64(0))
				bwDownGauge.Set(float64(0))
				netRateUpGauge.Set(float64(0))
				netRateDownGauge.Set(float64(0))
				vpnRateUpGauge.Set(float64(0))
				vpnRateDownGauge.Set(float64(0))
			} else {
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
				log.Print(err)
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
				log.Print(err)
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
