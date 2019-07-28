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
		t.Error("Expected no err, but go", err)
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
		t.Error("Expected no err, but go", err)
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
		t.Error("Expected no err, but go", err)
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
		t.Error("Expected no err, but go", err)
	}

	if cpum != 0 || cpub != 0 || sw != 0 || hdd != 0 || fanSpeed != 0 {
		t.Errorf("Expected 01 02 03 04 05, but got %v %v %v %v %v\n", cpum, cpub, sw, hdd, fanSpeed)
	}
}
