package main

import (
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
	version   string
	listen    string
	sessToken string
)

func init() {
	flag.StringVar(&mafreebox, "endpoint", "http://mafreebox.freebox.fr/", "Endpoint for freebox API")
	flag.StringVar(&version, "version", "v6", "freebox API version")
	flag.StringVar(&listen, "listen", ":10001", "prometheus metrics port")
}

func main() {
	flag.Parse()

	if !strings.HasSuffix(mafreebox, "/") {
		mafreebox = mafreebox + "/"
	}

	// myAuthInfo contains all auth data
	endpoint := mafreebox + "api/" + version + "/login/"
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
	}

	myPostRequest := newPostRequest()

	myLanRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/" + version + "/lan/browser/pub/",
		header: "X-Fbx-App-Auth",
	}

	mySystemRequest := &postRequest{
		method: "GET",
		url:    mafreebox + "api/" + version + "/system/",
		header: "X-Fbx-App-Auth",
	}

	// infinite loop to get all statistics
	go func() {
		for {
			// dsl metrics
			rateUp, rateDown, snrUp, snrDown, err := getDsl(myAuthInfo, myPostRequest)
			if err != nil {
				log.Fatalln(err)
			}
			rateUpGauge.Set(float64(rateUp))
			rateDownGauge.Set(float64(rateDown))
			snrUpGauge.Set(float64(snrUp))
			snrDownGauge.Set(float64(snrDown))

			// switch metrcis
			rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4, err := getSwitch(myAuthInfo, myPostRequest)
			if err != nil {
				log.Fatal(err)
			}
			rx1Gauge.Set(float64(rx1))
			tx1Gauge.Set(float64(tx1))
			rx2Gauge.Set(float64(rx2))
			tx2Gauge.Set(float64(tx2))
			rx3Gauge.Set(float64(rx3))
			tx3Gauge.Set(float64(tx3))
			rx4Gauge.Set(float64(rx4))
			tx4Gauge.Set(float64(tx4))

			// temps metrcis
			cpum, cpub, sw, hdd, fanSpeed, err := getTemp(myAuthInfo, myPostRequest)
			if err != nil {
				log.Fatal(err)
			}
			cpumGauge.Set(float64(cpum))
			cpubGauge.Set(float64(cpub))
			swGauge.Set(float64(sw))
			hddGauge.Set(float64(hdd))
			fanSpeedGauge.Set(float64(fanSpeed))

			// net metrics
			bwUp, bwDown, netRateUp, netRateDown, vpnRateUp, vpnRateDown, err := getNet(myAuthInfo, myPostRequest)
			if err != nil {
				log.Fatal(err)
			}
			bwUpGauge.Set(float64(bwUp))
			bwDownGauge.Set(float64(bwDown))
			netRateUpGauge.Set(float64(netRateUp))
			netRateDownGauge.Set(float64(netRateDown))
			vpnRateUpGauge.Set(float64(vpnRateUp))
			vpnRateDownGauge.Set(float64(vpnRateDown))

			// lan metrics
			lanAvailable, err := getLan(myAuthInfo, myLanRequest)
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range lanAvailable {
				if v.Reachable {
					lanReachableGauges.WithLabelValues(v.PrimaryName).Set(float64(1))
				} else {
					lanReachableGauges.WithLabelValues(v.PrimaryName).Set(float64(0))
				}
			}

			// fan metrics
			systemStats, err := getSystem(myAuthInfo, mySystemRequest)
			if err != nil {
				log.Fatal(err)
			}
			sensors := systemStats.Sensors
			fans := systemStats.Fans
			for _, v := range sensors {
				systemTempGauges.WithLabelValues(v.ID, v.Name).Set(float64(v.Value))
			}
			for _, v := range fans {
				systemFanGauges.WithLabelValues(v.ID, v.Name).Set(float64(v.Value))
			}

			time.Sleep(10 * time.Second)
		}
	}()

	// expose the registered metrics via HTTP OpenMetrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
