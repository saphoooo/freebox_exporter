package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// XXX: see https://dev.freebox.fr/sdk/os/ for API documentation
	// XXX: see https://prometheus.io/docs/practices/naming/ for metric names

	// connectionXdsl gauges
	connectionXdslStatusUptimeGauges = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_status_uptime_seconds_total",
	},
		[]string{
			"status",
			"protocol",
			"modulation",
		},
	)

	connectionXdslDownAttnGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_down_attn_decibels",
	})
	connectionXdslUpAttnGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_up_attn_decibels",
	})
	connectionXdslDownSnrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_down_snr_decibels",
	})
	connectionXdslUpSnrGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "freebox_connection_xdsl_up_snr_decibels",
	})

	connectionXdslErrorGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_errors_total",
			Help: "Error counts",
		},
		[]string{
			"direction", // up|down
			"name",      // crc|es|fec|hec
		},
	)

	connectionXdslGinpGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_ginp",
		},
		[]string{
			"direction", // up|down
			"name",      // enabled|rtx_(tx|c|uc)
		},
	)

	connectionXdslNitroGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_connection_xdsl_nitro",
		},
		[]string{
			"direction", // up|down
		},
	)

	// RRD dsl gauges
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

	// freeplug gauges
	freeplugRxRateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_rx_rate_bits",
		Help: "rx rate (from the freeplugs to the \"cco\" freeplug) (in bits/s) -1 if not available",
	},
		[]string{
			"id",
		},
	)
	freeplugTxRateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_tx_rate_bits",
		Help: "tx rate (from the \"cco\" freeplug to the freeplugs) (in bits/s) -1 if not available",
	},
		[]string{
			"id",
		},
	)
	freeplugHasNetworkGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "freebox_freeplug_has_network",
		Help: "is connected to the network",
	},
		[]string{
			"id",
		},
	)

	// RRD switch gauges
	// as switch database seems to be broken, this one is not used at this time
	/*
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
	*/

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

	// system temp gauges
	systemTempGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_temp_celsius",
			Help: "Temperature sensors reported by system (in °C)",
		},
		[]string{
			"name",
		},
	)

	// system fan gauges
	systemFanGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_fan_rpm",
			Help: "Fan speed reported by system (in RPM)",
		},
		[]string{
			"name",
		},
	)

	// system uptime gauges
	systemUptimeGauges = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "freebox_system_uptime_seconds_total",
		},
		[]string{
			"firmware_version",
		},
	)
)
