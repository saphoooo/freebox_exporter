package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSetFreeboxToken(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/login":
			myChall := &challenge{
				Success: true,
			}
			myChall.Result.Challenge = "foobar"
			result, _ := json.Marshal(myChall)
			fmt.Fprintln(w, string(result))
		case "/session":
			myToken := sessionToken{
				Success: true,
			}
			myToken.Result.SessionToken = "foobar"
			result, _ := json.Marshal(myToken)
			fmt.Fprintln(w, string(result))
		case "/granted/":
			myTrack := track{
				Success: true,
			}
			myTrack.Result.TrackID = 101
			result, _ := json.Marshal(myTrack)
			fmt.Fprintln(w, string(result))
		case "/granted/101":
			myGrant := grant{
				Success: true,
			}
			myGrant.Result.Status = "granted"
			result, _ := json.Marshal(myGrant)
			fmt.Fprintln(w, string(result))
		default:
			fmt.Fprintln(w, http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ai := &authInfo{}
	ai.myStore.location = "/tmp/token"
	ai.myAPI.login = ts.URL + "/login"
	ai.myAPI.loginSession = ts.URL + "/session"
	ai.myAPI.authz = ts.URL + "/granted/"
	ai.myReader = bufio.NewReader(strings.NewReader("\n"))

	var mySessionToken string

	token, err := setFreeboxToken(ai, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(ai.myStore.location)

	if token != "foobar" {
		t.Error("Expected foobar, but got", token)
	}

	os.Setenv("FREEBOX_TOKEN", "barfoo")
	defer os.Unsetenv("FREEBOX_TOKEN")

	token, err = setFreeboxToken(ai, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if token != "barfoo" {
		t.Error("Expected barfoo, but got", token)
	}
}

func TestGetDsl(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/good":
			myRRD := rrd{
				Success: true,
			}
			myRRD.Result.Data = []map[string]int64{
				{
					"rate_up":   12,
					"rate_down": 34,
					"snr_up":    56,
					"snr_down":  78,
				},
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/error":
			myRRD := rrd{
				Success:   true,
				ErrorCode: "insufficient_rights",
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/null":
			myRRD := rrd{
				Success: true,
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		}
	}))
	defer ts.Close()

	goodPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/good",
	}

	errorPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/error",
	}

	nullPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/null",
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	getDslResult, err := getDsl(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if getDslResult[0] != 12 || getDslResult[1] != 34 || getDslResult[2] != 56 || getDslResult[3] != 78 {
		t.Errorf("Expected 12 34 56 78, but got %v %v %v %v\n", getDslResult[0], getDslResult[1], getDslResult[2], getDslResult[3])
	}

	getDslResult, err = getDsl(ai, errorPR, &mySessionToken)
	if err.Error() != "Your app permissions does not allow accessing this API" {
		t.Error("Expected Your app permissions does not allow accessing this API, but go", err)
	}

	if len(getDslResult) != 0 {
		t.Error("Expected 0, but got", len(getDslResult))
	}

	getDslResult, err = getDsl(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if len(getDslResult) != 0 {
		t.Error("Expected 0, but got", len(getDslResult))
	}

}

func TestGetTemp(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/good":
			myRRD := rrd{
				Success: true,
			}
			myRRD.Result.Data = []map[string]int64{
				{
					"cpum":      01,
					"cpub":      02,
					"sw":        03,
					"hdd":       04,
					"fan_speed": 05,
				},
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/error":
			myRRD := rrd{
				Success:   true,
				ErrorCode: "denied_from_external_ip",
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/null":
			myRRD := rrd{
				Success: true,
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		}
	}))
	defer ts.Close()

	goodPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/good",
	}

	errorPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/error",
	}

	nullPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/null",
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	getTempResult, err := getTemp(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if getTempResult[0] != 01 || getTempResult[1] != 02 || getTempResult[2] != 03 || getTempResult[3] != 04 || getTempResult[4] != 05 {
		t.Errorf("Expected 01 02 03 04 05, but got %v %v %v %v %v\n", getTempResult[0], getTempResult[1], getTempResult[2], getTempResult[3], getTempResult[4])
	}

	getTempResult, err = getTemp(ai, errorPR, &mySessionToken)
	if err.Error() != "You are trying to get an app_token from a remote IP" {
		t.Error("Expected You are trying to get an app_token from a remote IP, but go", err)
	}

	if len(getTempResult) != 0 {
		t.Error("Expected 0, but got", len(getTempResult))
	}

	getTempResult, err = getTemp(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if len(getTempResult) != 0 {
		t.Error("Expected 0, but got", len(getTempResult))
	}

}

func TestGetNet(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/good":
			myRRD := rrd{
				Success: true,
			}
			myRRD.Result.Data = []map[string]int64{
				{
					"bw_up":         12500000000,
					"bw_down":       12500000000,
					"rate_up":       12500000000,
					"rate_down":     12500000000,
					"vpn_rate_up":   12500000000,
					"vpn_rate_down": 12500000000,
				},
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/error":
			myRRD := rrd{
				Success:   true,
				ErrorCode: "new_apps_denied",
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/null":
			myRRD := rrd{
				Success: true,
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		}
	}))
	defer ts.Close()

	goodPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/good",
	}

	errorPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/error",
	}

	nullPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/null",
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	getNetResult, err := getNet(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but go", err)
	}

	if getNetResult[0] != 12500000000 || getNetResult[1] != 12500000000 || getNetResult[2] != 12500000000 || getNetResult[3] != 12500000000 || getNetResult[4] != 12500000000 || getNetResult[5] != 12500000000 {
		t.Errorf("Expected 01 02 03 04 05 06, but got %v %v %v %v %v %v\n", getNetResult[0], getNetResult[1], getNetResult[2], getNetResult[3], getNetResult[4], getNetResult[5])
	}

	getNetResult, err = getNet(ai, errorPR, &mySessionToken)
	if err.Error() != "New application token request has been disabled" {
		t.Error("Expected New application token request has been disabled, but got", err)
	}

	if len(getNetResult) != 0 {
		t.Error("Expected 0, but got", len(getNetResult))
	}

	getNetResult, err = getNet(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if len(getNetResult) != 0 {
		t.Error("Expected 0, but got", len(getNetResult))
	}

}

func TestGetSwitch(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/good":
			myRRD := rrd{
				Success: true,
			}
			myRRD.Result.Data = []map[string]int64{
				{
					"rx_1": 01,
					"tx_1": 11,
					"rx_2": 02,
					"tx_2": 12,
					"rx_3": 03,
					"tx_3": 13,
					"rx_4": 04,
					"tx_4": 14,
				},
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/error":
			myRRD := rrd{
				Success:   true,
				ErrorCode: "apps_denied",
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		case "/null":
			myRRD := rrd{
				Success: true,
			}
			result, _ := json.Marshal(myRRD)
			fmt.Fprintln(w, string(result))
		}
	}))
	defer ts.Close()

	goodPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/good",
	}

	errorPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/error",
	}

	nullPR := &postRequest{
		method: "POST",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/null",
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	getSwitchResult, err := getSwitch(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if getSwitchResult[0] != 01 || getSwitchResult[1] != 11 || getSwitchResult[2] != 02 || getSwitchResult[3] != 12 || getSwitchResult[4] != 03 || getSwitchResult[5] != 13 || getSwitchResult[6] != 04 || getSwitchResult[7] != 14 {
		t.Errorf("Expected 01 11 02 12 03 13 04 14, but got %v %v %v %v %v %v %v %v\n", getSwitchResult[0], getSwitchResult[1], getSwitchResult[2], getSwitchResult[3], getSwitchResult[4], getSwitchResult[5], getSwitchResult[6], getSwitchResult[7])
	}

	getSwitchResult, err = getSwitch(ai, errorPR, &mySessionToken)
	if err.Error() != "API access from apps has been disabled" {
		t.Error("Expected API access from apps has been disabled, but got", err)
	}

	if len(getSwitchResult) != 0 {
		t.Error("Expected 0, but got", len(getSwitchResult))
	}

	getSwitchResult, err = getSwitch(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if len(getSwitchResult) != 0 {
		t.Error("Expected 0, but got", len(getSwitchResult))
	}

}

func TestGetLan(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/good":
			myLan := lan{
				Success: true,
			}
			myLan.Result = []lanHost{
				{
					Reachable:   true,
					PrimaryName: "Reachable host",
				},
				{
					Reachable:   false,
					PrimaryName: "Unreachable host",
				},
			}
			result, _ := json.Marshal(myLan)
			fmt.Fprintln(w, string(result))
		case "/error":
			myLan := lan{
				Success:   true,
				ErrorCode: "ratelimited",
			}
			result, _ := json.Marshal(myLan)
			fmt.Fprintln(w, string(result))
		}
	}))
	defer ts.Close()

	goodPR := &postRequest{
		method: "GET",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/good",
	}

	errorPR := &postRequest{
		method: "GET",
		header: "X-Fbx-App-Auth",
		url:    ts.URL + "/error",
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	lanAvailable, err := getLan(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	for _, v := range lanAvailable {
		if v.Reachable && v.PrimaryName != "Reachable host" {
			t.Errorf("Expected Reachable: true, Host: Reachable host, but go Reachable: %v, Host: %v", v.Reachable, v.PrimaryName)
		}

		if !v.Reachable && v.PrimaryName != "Unreachable host" {
			t.Errorf("Expected Reachable: false, Host: Unreachable host, but go Reachable: %v, Host: %v", !v.Reachable, v.PrimaryName)
		}
	}

	lanAvailable, err = getLan(ai, errorPR, &mySessionToken)
	if err.Error() != "Too many auth error have been made from your IP" {
		t.Error("Expected Too many auth error have been made from your IP, but got", err)
	}

}

func TestGetSystem(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mySys := system{
			Success: true,
		}
		mySys.Result.FanRPM = 666
		mySys.Result.TempCpub = 81
		mySys.Result.TempCpum = 89
		mySys.Result.TempHDD = 30
		mySys.Result.TempSW = 54

		/*
			mySys.Result {
				FanRPM:   666,
				TempCpub: 81,
				TempCpum: 89,
				TempHDD:  30,
				TempSW:   54,
			}
		*/
		result, _ := json.Marshal(mySys)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	pr := &postRequest{
		method: "GET",
		header: "X-Fbx-App-Auth",
		url:    ts.URL,
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	systemStats, err := getSystem(ai, pr, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if systemStats.Result.FanRPM != 666 {
		t.Error("Expected 666, but got", systemStats.Result.FanRPM)
	}

	if systemStats.Result.TempCpub != 81 {
		t.Error("Expected 81, but got", systemStats.Result.TempCpub)
	}

	if systemStats.Result.TempCpum != 89 {
		t.Error("Expected 89, but got", systemStats.Result.TempCpum)
	}

	if systemStats.Result.TempHDD != 30 {
		t.Error("Expected 30, but got", systemStats.Result.TempHDD)
	}

	if systemStats.Result.TempSW != 54 {
		t.Error("Expected 54, but got", systemStats.Result.TempSW)
	}

}

func TestGetWifi(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myWifi := wifi{
			Success: true,
		}
		myAP := wifiAccessPoint{
			Name: "AP1",
			ID:   0,
		}
		myWifi.Result = []wifiAccessPoint{myAP}

		result, _ := json.Marshal(myWifi)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	pr := &postRequest{
		method: "GET",
		header: "X-Fbx-App-Auth",
		url:    ts.URL,
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	wifiStats, err := getWifi(ai, pr, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if wifiStats.Result[0].Name != "AP1" {
		t.Error("Expected AP1, but got", wifiStats.Result[0].Name)
	}

	if wifiStats.Result[0].ID != 0 {
		t.Error("Expected 0, but got", wifiStats.Result[0].ID)
	}

}

func TestGetWifiStations(t *testing.T) {
	os.Setenv("FREEBOX_TOKEN", "IOI")
	defer os.Unsetenv("FREEBOX_TOKEN")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myWifiStations := wifiStations{
			Success: true,
		}

		myStation := wifiStation{
			Hostname:           "station_host",
			MAC:                "AA:BB:CC:DD:EE:FF",
			State:              "authorized",
			Inactive:           60,
			RXBytes:            500,
			TXBytes:            10000,
			ConnectionDuration: 600,
			TXRate:             20,
			RXRate:             5,
			Signal:             -20,
		}
		myWifiStations.Result = []wifiStation{myStation}

		result, _ := json.Marshal(myWifiStations)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	pr := &postRequest{
		method: "GET",
		header: "X-Fbx-App-Auth",
		url:    ts.URL,
	}

	ai := &authInfo{}
	mySessionToken := "foobar"

	wifiStationsStats, err := getWifiStations(ai, pr, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if wifiStationsStats.Result[0].Hostname != "station_host" {
		t.Error("Expected station_host, but got", wifiStationsStats.Result[0].Hostname)
	}

	if wifiStationsStats.Result[0].MAC != "AA:BB:CC:DD:EE:FF" {
		t.Error("Expected AA:BB:CC:DD:EE:FF, but got", wifiStationsStats.Result[0].MAC)
	}

	if wifiStationsStats.Result[0].State != "authorized" {
		t.Error("Expected authorized, but got", wifiStationsStats.Result[0].State)
	}

	if wifiStationsStats.Result[0].Inactive != 60 {
		t.Error("Expected 60, but got", wifiStationsStats.Result[0].Inactive)
	}

	if wifiStationsStats.Result[0].RXBytes != 500 {
		t.Error("Expected 500, but got", wifiStationsStats.Result[0].RXBytes)
	}

	if wifiStationsStats.Result[0].TXBytes != 10000 {
		t.Error("Expected 10000, but got", wifiStationsStats.Result[0].TXBytes)
	}

	if wifiStationsStats.Result[0].ConnectionDuration != 600 {
		t.Error("Expected 600, but got", wifiStationsStats.Result[0].ConnectionDuration)
	}

	if wifiStationsStats.Result[0].TXRate != 20 {
		t.Error("Expected 20, but got", wifiStationsStats.Result[0].TXRate)
	}

	if wifiStationsStats.Result[0].RXRate != 5 {
		t.Error("Expected 5, but got", wifiStationsStats.Result[0].RXRate)
	}

	if wifiStationsStats.Result[0].Signal != -20 {
		t.Error("Expected -20, but got", wifiStationsStats.Result[0].Signal)
	}

}

func Test_getNet(t *testing.T) {
	type args struct {
		authInf       *authInfo
		pr            *postRequest
		xSessionToken *string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNet(tt.args.authInf, tt.args.pr, tt.args.xSessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNet() = %v, want %v", got, tt.want)
			}
		})
	}
}
