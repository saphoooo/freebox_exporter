package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
			myRRD.Result.Data = []map[string]int{
				map[string]int{
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

	r1, r2, s1, s2, err := getDsl(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if r1 != 12 || r2 != 34 || s1 != 56 || s2 != 78 {
		t.Errorf("Expected 12 34 56 78, but got %v %v %v %v\n", r1, r2, s1, s2)
	}

	r1, r2, s1, s2, err = getDsl(ai, errorPR, &mySessionToken)
	if err.Error() != "Your app permissions does not allow accessing this API" {
		t.Error("Expected Your app permissions does not allow accessing this API, but go", err)
	}

	r1, r2, s1, s2, err = getDsl(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if r1 != 0 || r2 != 0 || s1 != 0 || s2 != 0 {
		t.Errorf("Expected 12 34 56 78, but got %v %v %v %v\n", r1, r2, s1, s2)
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
			myRRD.Result.Data = []map[string]int{
				map[string]int{
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

	cpum, cpub, sw, hdd, fanSpeed, err := getTemp(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if cpum != 01 || cpub != 02 || sw != 03 || hdd != 04 || fanSpeed != 05 {
		t.Errorf("Expected 01 02 03 04 05, but got %v %v %v %v %v\n", cpum, cpub, sw, hdd, fanSpeed)
	}

	cpum, cpub, sw, hdd, fanSpeed, err = getTemp(ai, errorPR, &mySessionToken)
	if err.Error() != "You are trying to get an app_token from a remote IP" {
		t.Error("Expected You are trying to get an app_token from a remote IP, but go", err)
	}

	cpum, cpub, sw, hdd, fanSpeed, err = getTemp(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if cpum != 0 || cpub != 0 || sw != 0 || hdd != 0 || fanSpeed != 0 {
		t.Errorf("Expected 01 02 03 04 05, but got %v %v %v %v %v\n", cpum, cpub, sw, hdd, fanSpeed)
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
			myRRD.Result.Data = []map[string]int{
				map[string]int{
					"bw_up":         01,
					"bw_down":       02,
					"rate_up":       03,
					"rate_down":     04,
					"vpn_rate_up":   05,
					"vpn_rate_down": 06,
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

	bwUP, bwDown, rUP, rDown, vpnUP, vpnDown, err := getNet(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but go", err)
	}

	if bwUP != 01 || bwDown != 02 || rUP != 03 || rDown != 04 || vpnUP != 05 || vpnDown != 06 {
		t.Errorf("Expected 01 02 03 04 05 06, but got %v %v %v %v %v %v\n", bwUP, bwDown, rUP, rDown, vpnUP, vpnDown)
	}

	bwUP, bwDown, rUP, rDown, vpnUP, vpnDown, err = getNet(ai, errorPR, &mySessionToken)
	if err.Error() != "New application token request has been disabled" {
		t.Error("Expected New application token request has been disabled, but got", err)
	}

	bwUP, bwDown, rUP, rDown, vpnUP, vpnDown, err = getNet(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if bwUP != 0 || bwDown != 0 || rUP != 0 || rDown != 0 || vpnUP != 0 || vpnDown != 0 {
		t.Errorf("Expected 01 02 03 04 05 06, but got %v %v %v %v %v %v\n", bwUP, bwDown, rUP, rDown, vpnUP, vpnDown)
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
			myRRD.Result.Data = []map[string]int{
				map[string]int{
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

	rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4, err := getSwitch(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if rx1 != 01 || tx1 != 11 || rx2 != 02 || tx2 != 12 || rx3 != 03 || tx3 != 13 || rx4 != 04 || tx4 != 14 {
		t.Errorf("Expected 01 11 02 12 03 13 04 14, but got %v %v %v %v %v %v %v %v\n", rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4)
	}

	rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4, err = getSwitch(ai, errorPR, &mySessionToken)
	if err.Error() != "API access from apps has been disabled" {
		t.Error("Expected API access from apps has been disabled, but got", err)
	}

	rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4, err = getSwitch(ai, nullPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if rx1 != 0 || tx1 != 0 || rx2 != 0 || tx2 != 0 || rx3 != 0 || tx3 != 0 || rx4 != 0 || tx4 != 0 {
		t.Errorf("Expected 0 0 0 0 0 0 0 0, but got %v %v %v %v %v %v %v %v\n", rx1, tx1, rx2, tx2, rx3, tx3, rx4, tx4)
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
				lanHost{
					Reachable:   true,
					PrimaryName: "Reachable host",
				},
				lanHost{
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
		switch r.RequestURI {
		case "/good":
			mySys := system{
				Success: true,
			}
			mySys.Result = systemR{
				Sensors: []idNameValue{
					idNameValue{
						ID:    "sensor01",
						Name:  "Sensor",
						Value: 01,
					},
				},
				Fans: []idNameValue{
					idNameValue{
						ID:    "fan01",
						Name:  "Fan",
						Value: 02,
					},
				},
			}
			result, _ := json.Marshal(mySys)
			fmt.Fprintln(w, string(result))
		case "/error":
			mySys := system{
				Success:   true,
				ErrorCode: "invalid_token",
			}
			result, _ := json.Marshal(mySys)
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

	systemStats, err := getSystem(ai, goodPR, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	sensors := systemStats.Sensors
	fans := systemStats.Fans
	for _, v := range sensors {
		if v.ID != "sensor01" || v.Name != "Sensor" || v.Value != 01 {
			t.Errorf("Expected sensor01 Sensor 01, but got %s %s %d", v.ID, v.Name, v.Value)
		}
	}
	for _, v := range fans {
		if v.ID != "fan01" || v.Name != "Fan" || v.Value != 02 {
			t.Errorf("Expected fan01 Fan 02, but got %s %s %d", v.ID, v.Name, v.Value)
		}
	}

	systemStats, err = getSystem(ai, errorPR, &mySessionToken)
	if err.Error() != "The app token you are trying to use is invalid or has been revoked" {
		t.Error("Expected The app token you are trying to use is invalid or has been revoked, but got", err)
	}
}
