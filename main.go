package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mafreebox    string
	version      string
	listen       string
	trackID      *track
	granted      *grant
	challenged   *challenge
	token        *sessionToken
	rrdTest      *rrd
	promExporter = app{
		AppID:      "fr.freebox.exporter",
		AppName:    "prom_exporter",
		AppVersion: "0.1",
		DeviceName: "laptop",
	}
	sessToken string
)

func init() {
	flag.StringVar(&mafreebox, "endpoint", "http://mafreebox.freebox.fr/", "Endpoint for freebox API")
	flag.StringVar(&version, "version", "v6", "freebox API version")
	flag.StringVar(&listen, "listen", ":10001", "prometheus metrics port")
}

func main() {
	flag.Parse()

	// dsl gauge
	rateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_rate_up",
		Help: "Available upload bandwidth (in byte/s)",
	})
	rateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_rate_down",
		Help: "Available download bandwidth (in byte/s)",
	})
	snrUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_up",
		Help: "Upload signal/noise ratio (in 1/10 dB)",
	})
	snrDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_down",
		Help: "Download signal/noise ratio (in 1/10 dB)",
	})

	// switch gauges
	rx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx1",
		Help: "Receive rate on port 1 (in byte/s)",
	})
	tx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx1",
		Help: "Transmit on port 1 (in byte/s)",
	})
	rx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx2",
		Help: "Receive rate on port 2 (in byte/s)",
	})
	tx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx2",
		Help: "Transmit on port 2 (in byte/s)",
	})
	rx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx3",
		Help: "Receive rate on port 3 (in byte/s)",
	})
	tx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx3",
		Help: "Transmit on port 3 (in byte/s)",
	})
	rx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx4",
		Help: "Receive rate on port 4 (in byte/s)",
	})
	tx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx4",
		Help: "Transmit on port 4 (in byte/s)",
	})

	// temp gauges
	cpumGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpum",
		Help: "Temperature cpum (in °C)",
	})
	cpubGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpub",
		Help: "Temperature cpub (in °C)",
	})
	swGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_sw",
		Help: "Temperature sw (in °C)",
	})
	hddGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_hdd",
		Help: "Temperature hdd (in °C)",
	})
	fanSpeedGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_fan_speed",
		Help: "Fan rpm",
	})

	// net gauges
	bwUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_up",
		Help: "Upload available bandwidth (in byte/s)",
	})
	bwDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_down",
		Help: "Download available bandwidth (in byte/s)",
	})
	netRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_rate_up",
		Help: "Upload rate (in byte/s)",
	})
	netRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_rate_down",
		Help: "Download rate (in byte/s)",
	})
	vpnRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_rate_up",
		Help: "Vpn client upload rate (in byte/s)",
	})
	vpnRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_rate_down",
		Help: "Vpn client download rate (in byte/s)",
	})
	// infinite loop to get all statistics
	go func() {
		for {
			// dsl metrics
			rateUp, rateDown, snrUp, snrDown := getDsl()
			rateUpGauge.Set(float64(rateUp))
			rateDownGauge.Set(float64(rateDown))
			snrUpGauge.Set(float64(snrUp))
			snrDownGauge.Set(float64(snrDown))
			//fmt.Printf("rate_up: %d\nrate_down: %d\nsnr_up: %d\nsnr_down: %d\n", rateUp, rateDown, snrUp, snrDown)

			// switch metrcis
			rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4 := getSwitch()
			rx1Gauge.Set(float64(rx1))
			tx1Gauge.Set(float64(tx1))
			rx2Gauge.Set(float64(rx2))
			tx2Gauge.Set(float64(tx2))
			rx3Gauge.Set(float64(rx3))
			tx3Gauge.Set(float64(tx3))
			rx4Gauge.Set(float64(rx4))
			tx4Gauge.Set(float64(tx4))
			//fmt.Printf("rx1: %d\ntx1: %d\nrx2: %d\ntx2: %d\nrx3: %d\ntx3: %d\nrx4: %d\ntx4: %d\n", rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4)

			// temps metrcis
			cpum, cpub, sw, hdd, fanSpeed := getTemp()
			cpumGauge.Set(float64(cpum))
			cpubGauge.Set(float64(cpub))
			swGauge.Set(float64(sw))
			hddGauge.Set(float64(hdd))
			fanSpeedGauge.Set(float64(fanSpeed))
			//fmt.Printf("cpum: %d\ncpub: %d\nsw: %d\nhdd: %d\nfan_speed: %d\n", cpum, cpub, sw, hdd, fanSpeed)

			// net metrics
			bwUp, bwDown, netRateUp, netRateDown, vpnRateUp, vpnRateDown := getNet()
			bwUpGauge.Set(float64(bwUp))
			bwDownGauge.Set(float64(bwDown))
			netRateUpGauge.Set(float64(netRateUp))
			netRateDownGauge.Set(float64(netRateDown))
			vpnRateUpGauge.Set(float64(vpnRateUp))
			vpnRateDownGauge.Set(float64(vpnRateDown))
			//fmt.Printf("bw_up: %d\nbw_down: %d\nrate_up: %d\nrate_down: %d\nvpn_rate_up: %d\nvpn_rate_down: %d\n", bwUp, bwDown, rateUp, rateDown, vpnRateUp, vpnRateDown)

			time.Sleep(10 * time.Second)
		}
	}()

	// expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listen, nil))
}
