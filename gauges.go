package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RRD dsl gauge
	rateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_up_bytes",
		Help: "Available upload bandwidth (in byte/s)",
	})
	rateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_down_bytes",
		Help: "Available download bandwidth (in byte/s)",
	})
	snrUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_up_decibel",
		Help: "Upload signal/noise ratio (in 1/10 dB)",
	})
	snrDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_dsl_snr_down_decibel",
		Help: "Download signal/noise ratio (in 1/10 dB)",
	})

	// RRD switch gauges
	rx1Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx1_bytes",
		Help: "Receive rate on port 1 (in byte/s)",
	})
	tx1Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx1_bytes",
		Help: "Transmit on port 1 (in byte/s)",
	})
	rx2Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx2_bytes",
		Help: "Receive rate on port 2 (in byte/s)",
	})
	tx2Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx2_bytes",
		Help: "Transmit on port 2 (in byte/s)",
	})
	rx3Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx3_bytes",
		Help: "Receive rate on port 3 (in byte/s)",
	})
	tx3Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx3_bytes",
		Help: "Transmit on port 3 (in byte/s)",
	})
	rx4Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_rx4_bytes",
		Help: "Receive rate on port 4 (in byte/s)",
	})
	tx4Gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_switch_tx4_bytes",
		Help: "Transmit on port 4 (in byte/s)",
	})

	// RRD temp gauges
	// these ones look broken, use system temp gauges instead
	cpumGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpum_celsius",
		Help: "Temperature cpum (in °C)",
	})
	cpubGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_cpub_celsius",
		Help: "Temperature cpub (in °C)",
	})
	swGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_sw_celsius",
		Help: "Temperature sw (in °C)",
	})
	hddGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_hdd_celsius",
		Help: "Temperature hdd (in °C)",
	})
	fanSpeedGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_temp_fan_speed_rpm",
		Help: "Fan rpm",
	})

	// RRD net gauges
	bwUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_up_bytes",
		Help: "Upload available bandwidth (in byte/s)",
	})
	bwDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_bw_down_bytes",
		Help: "Download available bandwidth (in byte/s)",
	})
	netRateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_up_bytes",
		Help: "Upload rate (in byte/s)",
	})
	netRateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_down_bytes",
		Help: "Download rate (in byte/s)",
	})
	vpnRateUpGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_up_bytes",
		Help: "Vpn client upload rate (in byte/s)",
	})
	vpnRateDownGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_net_vpn_down_bytes",
		Help: "Vpn client download rate (in byte/s)",
	})

	// lan gauges
	lanReachableGauges = promauto.NewGaugeVec(
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
	systemTempGauges = promauto.NewGaugeVec(
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
	systemFanGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_fan_rpm",
			Help: "Fan speed reported by system (in rpm)",
		},
		[]string{
			"id",
			"name",
		},
	)
)
