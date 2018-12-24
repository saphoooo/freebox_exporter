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
		Help: "DSL rate up",
	})
	rateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_rate_down",
		Help: "DSL rate down",
	})
	snrUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_up",
		Help: "DSL snr up",
	})
	snrDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_down",
		Help: "DSL snr down",
	})

	// switch gauges
	rx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx1",
		Help: "SWITCH rx1",
	})
	tx1Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx1",
		Help: "SWITCH tx1",
	})
	rx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx2",
		Help: "SWITCH rx2",
	})
	tx2Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx2",
		Help: "SWITCH tx2",
	})
	rx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx3",
		Help: "SWITCH rx3",
	})
	tx3Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx3",
		Help: "SWITCH tx3",
	})
	rx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx4",
		Help: "SWITCH rx4",
	})
	tx4Gauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx4",
		Help: "SWITCH tx4",
	})

	// temp gauges
	cpumGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpum",
		Help: "TEMP cpum",
	})
	cpubGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpub",
		Help: "TEMP cpub",
	})
	swGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_sw",
		Help: "TEMP sw",
	})
	hddGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_hdd",
		Help: "TEMP hdd",
	})
	fanSpeedGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_fan_speed",
		Help: "TEMP fan speed",
	})

	// net gauges
	bwUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_up",
		Help: "NET bw up",
	})
	bwDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_down",
		Help: "NET bw down",
	})
	netRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_rate_up",
		Help: "NET rate up",
	})
	netRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_rate_down",
		Help: "NET rate down",
	})
	vpnRateUpGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_rate_up",
		Help: "NET vpn rate up",
	})
	vpnRateDownGauge := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_rate_down",
		Help: "NET vpn rate down",
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
