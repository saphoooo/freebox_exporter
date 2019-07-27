package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mafreebox  string
	version    string
	listen     string
	trackID    *track
	granted    *grant
	challenged *challenge
	token      *sessionToken
	rrdTest    *rrd
	lanResp    *lan
	systemResp *system
	sessToken  string
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

	// RRD dsl gauge
	rateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_up_bytes",
		Help: "Available upload bandwidth (in byte/s)",
	})
	rateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_down_bytes",
		Help: "Available download bandwidth (in byte/s)",
	})
	snrUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_up_decibel",
		Help: "Upload signal/noise ratio (in 1/10 dB)",
	})
	snrDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_down_decibel",
		Help: "Download signal/noise ratio (in 1/10 dB)",
	})

	// RRD switch gauges
	rx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx1_bytes",
		Help: "Receive rate on port 1 (in byte/s)",
	})
	tx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx1_bytes",
		Help: "Transmit on port 1 (in byte/s)",
	})
	rx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx2_bytes",
		Help: "Receive rate on port 2 (in byte/s)",
	})
	tx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx2_bytes",
		Help: "Transmit on port 2 (in byte/s)",
	})
	rx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx3_bytes",
		Help: "Receive rate on port 3 (in byte/s)",
	})
	tx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx3_bytes",
		Help: "Transmit on port 3 (in byte/s)",
	})
	rx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx4_bytes",
		Help: "Receive rate on port 4 (in byte/s)",
	})
	tx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx4_bytes",
		Help: "Transmit on port 4 (in byte/s)",
	})

	// RRD temp gauges
	// these ones look broken, use system temp gauges instead
	cpumGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpum_celsius",
		Help: "Temperature cpum (in °C)",
	})
	cpubGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpub_celsius",
		Help: "Temperature cpub (in °C)",
	})
	swGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_sw_celsius",
		Help: "Temperature sw (in °C)",
	})
	hddGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_hdd_celsius",
		Help: "Temperature hdd (in °C)",
	})
	fanSpeedGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_fan_speed_rpm",
		Help: "Fan rpm",
	})

	// RRD net gauges
	bwUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_up_bytes",
		Help: "Upload available bandwidth (in byte/s)",
	})
	bwDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_down_bytes",
		Help: "Download available bandwidth (in byte/s)",
	})
	netRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_up_bytes",
		Help: "Upload rate (in byte/s)",
	})
	netRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_down_bytes",
		Help: "Download rate (in byte/s)",
	})
	vpnRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_up_bytes",
		Help: "Vpn client upload rate (in byte/s)",
	})
	vpnRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_down_bytes",
		Help: "Vpn client download rate (in byte/s)",
	})

	// lan gauges
	lanReachableGauges := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_lan_reachable",
			Help: "Hosts reachable on LAN",
		},
		[]string{
			// Hostname
			"name",
		},
	)
	// SYSTEM temp gauges
	systemTempGauges := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_temp_celsius",
			Help: "Temp sensors reported by system (in °C)",
		},
		[]string{
			"id",
			"name",
		},
	)

	// SYSTEM fan gauges
	systemFanGauges := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_fan_rpm",
			Help: "Fan speed reported by system (in rpm)",
		},
		[]string{
			"id",
			"name",
		},
	)

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
			systemStats := getSystem(myAuthInfo)
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
