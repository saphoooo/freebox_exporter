package main

import (
	"fmt"
	"time"
)

func main() {
	// infinite loop to get all statistics
	for {
		rateUp, rateDown, snrUp, snrDown := getDsl()
		fmt.Printf("rate_up: %d\nrate_down: %d\nsnr_up: %d\nsnr_down: %d\n", rateUp, rateDown, snrUp, snrDown)
		rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4 := getSwitch()
		fmt.Printf("rx1: %d\ntx1: %d\nrx2: %d\ntx2: %d\nrx3: %d\ntx3: %d\nrx4: %d\ntx4: %d\n", rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4)
		cpum, cpub, sw, hdd, fanSpeed := getTemp()
		fmt.Printf("cpum: %d\ncpub: %d\nsw: %d\nhdd: %d\nfan_speed: %d\n", cpum, cpub, sw, hdd, fanSpeed)
		bwUp, bwDown, rateUp, rateDown, vpnRateUp, vpnRateDown := getNet()
		fmt.Printf("bw_up: %d\nbw_down: %d\nrate_up: %d\nrate_down: %d\nvpn_rate_up: %d\nvpn_rate_down: %d\n", bwUp, bwDown, rateUp, rateDown, vpnRateUp, vpnRateDown)
		time.Sleep(10 * time.Second)
	}
}
