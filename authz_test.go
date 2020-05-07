package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRetreiveToken(t *testing.T) {
	ai := &authInfo{
		myStore: store{location: "/tmp/token"},
	}

	_, err := retreiveToken(ai)
	if err.Error() != "stat /tmp/token: no such file or directory" {
		t.Error("Expected bla, but got", err)
	}

	ioutil.WriteFile(ai.myStore.location, []byte("IOI"), 0600)
	defer os.Remove(ai.myStore.location)

	token, err := retreiveToken(ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	newToken := os.Getenv("FREEBOX_TOKEN")

	if newToken != "IOI" {
		t.Error("Expected IOI, but got", newToken)
	}

	if token != "IOI" {
		t.Error("Expected IOI, but got", newToken)
	}

	os.Unsetenv("FREEBOX_TOKEN")
}

func TestStoreToken(t *testing.T) {
	var token string

	ai := &authInfo{}
	token = "IOI"
	err := storeToken(token, ai)
	if err.Error() != "open : no such file or directory" {
		t.Error("Expected open : no such file or directory, but got", err)
	}

	ai.myStore.location = "/tmp/token"
	err = storeToken(token, ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(ai.myStore.location)

	token = os.Getenv("FREEBOX_TOKEN")
	if token != "IOI" {
		t.Error("Expected IOI, but got", token)
	}
	os.Unsetenv("FREEBOX_TOKEN")

	data, err := ioutil.ReadFile(ai.myStore.location)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if string(data) != "IOI" {
		t.Error("Expected IOI, but got", string(data))
	}

}

func TestGetTrackID(t *testing.T) {
	ai := &authInfo{
		myStore: store{location: "/tmp/token"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myTrack := track{
			Success: true,
		}
		myTrack.Result.AppToken = "IOI"
		myTrack.Result.TrackID = 101
		result, _ := json.Marshal(myTrack)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	ai.myAPI.authz = ts.URL
	trackID, err := getTrackID(ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(ai.myStore.location)
	defer os.Unsetenv("FREEBOX_TOKEN")

	if trackID.Result.TrackID != 101 {
		t.Error("Expected 101, but got", trackID.Result.TrackID)
	}

	// as getTrackID have no return value
	// the result of storeToken func is checked instead
	token := os.Getenv("FREEBOX_TOKEN")
	if token != "IOI" {
		t.Error("Expected IOI, but got", token)
	}
}

func TestGetGranted(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/unknown/":
			myTrack := track{
				Success: true,
			}
			myTrack.Result.TrackID = 101
			result, _ := json.Marshal(myTrack)
			fmt.Fprintln(w, string(result))
		case "/unknown/101":
			myGrant := grant{
				Success: true,
			}
			myGrant.Result.Status = "unknown"
			result, _ := json.Marshal(myGrant)
			fmt.Fprintln(w, string(result))
		case "/timeout/":
			myTrack := track{
				Success: true,
			}
			myTrack.Result.TrackID = 101
			result, _ := json.Marshal(myTrack)
			fmt.Fprintln(w, string(result))
		case "/timeout/101":
			myGrant := grant{
				Success: true,
			}
			myGrant.Result.Status = "timeout"
			result, _ := json.Marshal(myGrant)
			fmt.Fprintln(w, string(result))
		case "/denied/":
			myTrack := track{
				Success: true,
			}
			myTrack.Result.TrackID = 101
			result, _ := json.Marshal(myTrack)
			fmt.Fprintln(w, string(result))
		case "/denied/101":
			myGrant := grant{
				Success: true,
			}
			myGrant.Result.Status = "denied"
			result, _ := json.Marshal(myGrant)
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

	ai := authInfo{}
	ai.myAPI.authz = ts.URL + "/unknown/"
	ai.myStore.location = "/tmp/token"

	err := getGranted(&ai)
	if err.Error() != "the app_token is invalid or has been revoked" {
		t.Error("Expected the app_token is invalid or has been revoked, but got", err)
	}
	defer os.Remove(ai.myStore.location)
	defer os.Unsetenv("FREEBOX_TOKEN")

	ai.myAPI.authz = ts.URL + "/timeout/"
	err = getGranted(&ai)
	if err.Error() != "the user did not confirmed the authorization within the given time" {
		t.Error("Expected the user did not confirmed the authorization within the given time, but got", err)
	}

	ai.myAPI.authz = ts.URL + "/denied/"
	err = getGranted(&ai)
	if err.Error() != "the user denied the authorization request" {
		t.Error("Expected the user denied the authorization request, but got", err)
	}

	ai.myAPI.authz = ts.URL + "/granted/"
	err = getGranted(&ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
}

func TestGetChallenge(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myChall := &challenge{
			Success: true,
		}
		myChall.Result.Challenge = "foobar"
		result, _ := json.Marshal(myChall)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	ai := &authInfo{
		myAPI: api{
			login: ts.URL,
		},
	}

	challenged, err := getChallenge(ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if challenged.Success != true {
		t.Error("Expected true, but got", challenged.Success)
	}

	if challenged.Result.Challenge != "foobar" {
		t.Error("Expected foobar, but got", challenged.Result.Challenge)
	}
}

func TestHmacSha1(t *testing.T) {
	hmac := hmacSha1("IOI", "foobar")
	if hmac != "02fb876a39b64eddcfee3eaa69465cb3e8d53cde" {
		t.Error("Expected 02fb876a39b64eddcfee3eaa69465cb3e8d53cde, but got", hmac)
	}
}

func TestGetSession(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myToken := &sessionToken{
			Success: true,
		}
		myToken.Result.Challenge = "foobar"
		result, _ := json.Marshal(myToken)
		fmt.Fprintln(w, string(result))
	}))
	defer ts.Close()

	ai := &authInfo{
		myAPI: api{
			loginSession: ts.URL,
		},
	}

	token, err := getSession(ai, "")
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Unsetenv("FREEBOX_TOKEN")

	if token.Success != true {
		t.Error("Expected true, but got", token.Success)
	}

	if token.Result.Challenge != "foobar" {
		t.Error("Expected foobar, but got", token.Result.Challenge)
	}
}

func TestGetToken(t *testing.T) {

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

	ai := authInfo{}
	ai.myStore.location = "/tmp/token"
	ai.myAPI.login = ts.URL + "/login"
	ai.myAPI.loginSession = ts.URL + "/session"
	ai.myAPI.authz = ts.URL + "/granted/"
	ai.myReader = bufio.NewReader(strings.NewReader("\n"))

	var mySessionToken string

	// the first pass valide getToken without a token stored in a file
	tk, err := getToken(&ai, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(ai.myStore.location)
	defer os.Unsetenv("FREEBOX_TOKEN")

	if mySessionToken != "foobar" {
		t.Error("Expected foobar, but got", mySessionToken)
	}

	if tk != "foobar" {
		t.Error("Expected foobar, but got", tk)
	}

	// the second pass validate getToken with a token stored in a file:
	// the first pass creates a file at ai.myStore.location
	tk, err = getToken(&ai, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if mySessionToken != "foobar" {
		t.Error("Expected foobar, but got", mySessionToken)
	}

	if tk != "foobar" {
		t.Error("Expected foobar, but got", tk)
	}

}

func TestGetSessToken(t *testing.T) {

	myToken := &sessionToken{}

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
			myToken.Success = true
			myToken.Result.SessionToken = "foobar"
			result, _ := json.Marshal(myToken)
			fmt.Fprintln(w, string(result))
		case "/session2":
			myToken.Msg = "failed to get a session"
			myToken.Success = false
			result, _ := json.Marshal(myToken)
			fmt.Fprintln(w, string(result))
		default:
			fmt.Fprintln(w, http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ai := authInfo{}
	ai.myAPI.login = ts.URL + "/login"
	ai.myAPI.loginSession = ts.URL + "/session"
	var mySessionToken string

	st, err := getSessToken("token", &ai, &mySessionToken)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if st != "foobar" {
		t.Error("Expected foobar, but got", st)
	}

	ai.myAPI.loginSession = ts.URL + "/session2"

	_, err = getSessToken("token", &ai, &mySessionToken)
	if err.Error() != "failed to get a session" {
		t.Error("Expected but got failed to get a session, but got", err)
	}
}
