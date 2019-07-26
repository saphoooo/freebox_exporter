package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	apiErrors = map[string]error{
		"invalid_token":           errors.New("The app token you are trying to use is invalid or has been revoked"),
		"insufficient_rights":     errors.New("Your app permissions does not allow accessing this API"),
		"denied_from_external_ip": errors.New("You are trying to get an app_token from a remote IP"),
		"invalid_request":         errors.New("Your request is invalid"),
		"ratelimited":             errors.New("Too many auth error have been made from your IP"),
		"new_apps_denied":         errors.New("New application token request has been disabled"),
		"apps_denied":             errors.New("API access from apps has been disabled"),
		"internal_error":          errors.New("Internal error"),
	}
)

// setFreeboxToken ensure that there is an active token for a call
func setFreeboxToken(fb *freebox, st *store) string {
	token := os.Getenv("FREEBOX_TOKEN")
	if token == "" {
		sessToken = getToken(fb, st)
	}
	if sessToken == "" {
		sessToken = getSessToken(token)
	}
	return token
}

func newPostRequest() *postRequest {
	return &postRequest{
		method: "POST",
		url:    mafreebox + "api/" + version + "/rrd/",
		header: "X-Fbx-App-Auth",
	}
}

// getDsl get dsl statistics
func (d *database) getDsl(fb *freebox, st *store, pr *postRequest) (int, int, int, int, error) {
	freeboxToken := setFreeboxToken(fb, st)
	client := http.Client{}
	r, err := json.Marshal(d)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	req.Header.Add(pr.header, sessToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if resp.StatusCode == 404 {
		return 0, 0, 0, 0, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		sessToken = getSessToken(freeboxToken)
	case "invalid_token":
		return 0, 0, 0, 0, apiErrors["auth_required"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return 0, 0, 0, 0, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return 0, 0, 0, 0, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return 0, 0, 0, 0, apiErrors["invalid_request"]
	case "ratelimited":
		return 0, 0, 0, 0, apiErrors["ratelimited"]
	case "new_apps_denied":
		return 0, 0, 0, 0, apiErrors["new_apps_denied"]
	case "apps_denied":
		return 0, 0, 0, 0, apiErrors["apps_denied"]
	case "internal_error":
		return 0, 0, 0, 0, apiErrors["internal_error"]
	}
	if len(rrdTest.Result.Data) == 0 {
		return 0, 0, 0, 0, nil
	}
	return rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["snr_up"], rrdTest.Result.Data[0]["snr_down"], nil
}
