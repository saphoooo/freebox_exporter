package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
func setFreeboxToken(authInf *authInfo, xSessionToken *string) (string, error) {

	token := os.Getenv("FREEBOX_TOKEN")

	if token == "" {
		var err error
		*xSessionToken, err = getToken(authInf, xSessionToken)
		if err != nil {
			return "", err
		}
	}

	if *xSessionToken == "" {
		var err error
		*xSessionToken, err = getSessToken(token, authInf, xSessionToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	return token, nil

}

func newPostRequest() *postRequest {
	return &postRequest{
		method: "POST",
		url:    mafreebox + "api/" + version + "/rrd/",
		header: "X-Fbx-App-Auth",
	}
}

// getDsl get dsl statistics
func getDsl(authInf *authInfo, pr *postRequest, xSessionToken *string) (int, int, int, int, error) {
	d := &database{
		DB:        "dsl",
		Fields:    []string{"rate_up", "rate_down", "snr_up", "snr_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	req.Header.Add(pr.header, *xSessionToken)
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
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			log.Fatal(err)
		}
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

// getTemp get temp statistics
func getTemp(authInf *authInfo, pr *postRequest, xSessionToken *string) (int, int, int, int, int, error) {
	d := &database{
		DB:        "temp",
		Fields:    []string{"cpum", "cpub", "sw", "hdd", "fan_speed"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, fmt.Sprintf(pr.url), buf)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	if resp.StatusCode == 404 {
		return 0, 0, 0, 0, 0, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return 0, 0, 0, 0, 0, err
		}
	case "invalid_token":
		return 0, 0, 0, 0, 0, apiErrors["invalid_token"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return 0, 0, 0, 0, 0, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return 0, 0, 0, 0, 0, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return 0, 0, 0, 0, 0, apiErrors["invalid_request"]
	case "ratelimited":
		return 0, 0, 0, 0, 0, apiErrors["ratelimited"]
	case "new_apps_denied":
		return 0, 0, 0, 0, 0, apiErrors["new_apps_denied"]
	case "apps_denied":
		return 0, 0, 0, 0, 0, apiErrors["apps_denied"]
	case "internal_error":
		return 0, 0, 0, 0, 0, apiErrors["internal_error"]
	}
	if len(rrdTest.Result.Data) == 0 {
		return 0, 0, 0, 0, 0, nil
	}
	return rrdTest.Result.Data[0]["cpum"], rrdTest.Result.Data[0]["cpub"], rrdTest.Result.Data[0]["sw"], rrdTest.Result.Data[0]["hdd"], rrdTest.Result.Data[0]["fan_speed"], nil
}

// getNet get net statistics
func getNet(authInf *authInfo, pr *postRequest, xSessionToken *string) (int, int, int, int, int, int, error) {
	d := &database{
		DB:        "net",
		Fields:    []string{"bw_up", "bw_down", "rate_up", "rate_down", "vpn_rate_up", "vpn_rate_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	if resp.StatusCode == 404 {
		return 0, 0, 0, 0, 0, 0, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, err
		}
	case "invalid_token":
		return 0, 0, 0, 0, 0, 0, apiErrors["invalid_token"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return 0, 0, 0, 0, 0, 0, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return 0, 0, 0, 0, 0, 0, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return 0, 0, 0, 0, 0, 0, apiErrors["invalid_request"]
	case "ratelimited":
		return 0, 0, 0, 0, 0, 0, apiErrors["ratelimited"]
	case "new_apps_denied":
		return 0, 0, 0, 0, 0, 0, apiErrors["new_apps_denied"]
	case "apps_denied":
		return 0, 0, 0, 0, 0, 0, apiErrors["apps_denied"]
	case "internal_error":
		return 0, 0, 0, 0, 0, 0, apiErrors["internal_error"]
	}
	if len(rrdTest.Result.Data) == 0 {
		return 0, 0, 0, 0, 0, 0, nil
	}
	return rrdTest.Result.Data[0]["bw_up"], rrdTest.Result.Data[0]["bw_down"], rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["vpn_rate_up"], rrdTest.Result.Data[0]["vpn_rate_down"], nil
}

// getSwitch get switch statistics
func getSwitch(authInf *authInfo, pr *postRequest, xSessionToken *string) (int, int, int, int, int, int, int, int, error) {
	d := &database{
		DB:        "switch",
		Fields:    []string{"rx_1", "tx_1", "rx_2", "tx_2", "rx_3", "tx_3", "rx_4", "tx_4"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	if resp.StatusCode == 404 {
		return 0, 0, 0, 0, 0, 0, 0, 0, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return 0, 0, 0, 0, 0, 0, 0, 0, err
		}
	case "invalid_token":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["invalid_token"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["invalid_request"]
	case "ratelimited":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["ratelimited"]
	case "new_apps_denied":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["new_apps_denied"]
	case "apps_denied":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["apps_denied"]
	case "internal_error":
		return 0, 0, 0, 0, 0, 0, 0, 0, apiErrors["internal_error"]
	}
	if len(rrdTest.Result.Data) == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, nil
	}
	return rrdTest.Result.Data[0]["rx_1"], rrdTest.Result.Data[0]["tx_1"], rrdTest.Result.Data[0]["rx_2"], rrdTest.Result.Data[0]["tx_2"], rrdTest.Result.Data[0]["rx_3"], rrdTest.Result.Data[0]["tx_3"], rrdTest.Result.Data[0]["rx_4"], rrdTest.Result.Data[0]["tx_4"], nil
}

// getLan get lan statistics
func getLan(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]lanHost, error) {

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []lanHost{}, err
	}

	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return []lanHost{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []lanHost{}, err
	}
	if resp.StatusCode == 404 {
		return []lanHost{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []lanHost{}, err
	}

	lanResp := lan{}
	err = json.Unmarshal(body, &lanResp)
	switch lanResp.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []lanHost{}, err
		}
	case "invalid_token":
		return []lanHost{}, apiErrors["invalid_token"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return []lanHost{}, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return []lanHost{}, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return []lanHost{}, apiErrors["invalid_request"]
	case "ratelimited":
		return []lanHost{}, apiErrors["ratelimited"]
	case "new_apps_denied":
		return []lanHost{}, apiErrors["new_apps_denied"]
	case "apps_denied":
		return []lanHost{}, apiErrors["apps_denied"]
	case "internal_error":
		return []lanHost{}, apiErrors["internal_error"]
	}
	return lanResp.Result, nil
}

// getLan get lan statistics
func getSystem(authInf *authInfo, pr *postRequest, xSessionToken *string) (systemR, error) {
	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return systemR{}, err
	}

	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return systemR{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return systemR{}, err
	}
	if resp.StatusCode == 404 {
		return systemR{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return systemR{}, err
	}
	systemResp := system{}
	err = json.Unmarshal(body, &systemResp)
	switch systemResp.ErrorCode {
	case "auth_required":
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return systemR{}, err
		}
	case "invalid_token":
		return systemR{}, apiErrors["invalid_token"]
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		return systemR{}, apiErrors["insufficient_rights"]
	case "denied_from_external_ip":
		return systemR{}, apiErrors["denied_from_external_ip"]
	case "invalid_request":
		return systemR{}, apiErrors["invalid_request"]
	case "ratelimited":
		return systemR{}, apiErrors["ratelimited"]
	case "new_apps_denied":
		return systemR{}, apiErrors["new_apps_denied"]
	case "apps_denied":
		return systemR{}, apiErrors["apps_denied"]
	case "internal_error":
		return systemR{}, apiErrors["internal_error"]
	}
	return systemResp.Result, nil
}
