package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

	err := storeToken(token, ai)
	if err.Error() != "token should not be blank" {
		t.Error("Expected token should not be blank, but got", err)
	}

	token = "IOI"
	err = storeToken(token, ai)
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
	err := getTrackID(ai)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(ai.myStore.location)
	defer os.Unsetenv("FREEBOX_TOKEN")

	// as getTrackID have no return value
	// the result of storeToken func is checked instead
	token := os.Getenv("FREEBOX_TOKEN")
	if token != "IOI" {
		t.Error("Expected IOI, but got", token)
	}
}

func TestGetGranted(t *testing.T) {
	/*
			type grant struct {
			Success bool `json:"success"`
			Result  struct {
				Status    string `json:"status"`
				Challenge string `json:"challenge"`
			} `json:"result"`
		}
	*/

	/*
		app := &app{}
		fb := &freebox{}
		st := &store{
			location: "/tmp/token",
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			myGrant := &grant{
				Success: true,
			}
			myGrant.Result.Status = "unknown"
			myGrant.Result.Challenge = ""
			result, _ := json.Marshal(myGrant)
			fmt.Fprintln(w, string(result))
		}))
		defer ts.Close()
	*/
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

	err := getChallenge(ai)
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

	err := getSession(ai, "")
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if token.Success != true {
		t.Error("Expected true, but got", token.Success)
	}

	if token.Result.Challenge != "foobar" {
		t.Error("Expected foobar, but got", token.Result.Challenge)
	}
}
