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
		"db_error":                errors.New("Oops, The database you are trying to access doesn't seem to exist"),
	}
)

func (r *rrd) status() error {
	if apiErrors[r.ErrorCode] == nil {
		return errors.New("RRD: The API returns an unknown error_code: " + r.ErrorCode)
	}
	return apiErrors[r.ErrorCode]
}

func (l *lan) status() error {
	if apiErrors[l.ErrorCode] == nil {
		return errors.New("LAN: The API returns an unknown error_code: " + l.ErrorCode)
	}
	return apiErrors[l.ErrorCode]
}

// setFreeboxToken ensure that there is an active token for a call
func setFreeboxToken(authInf *authInfo, xSessionToken *string) (string, error) {

	token := os.Getenv("FREEBOX_TOKEN")

	if token == "" {
		var err error
		*xSessionToken, err = getToken(authInf, xSessionToken)
		if err != nil {
			return "", err
		}
		token = *xSessionToken
	}

	if *xSessionToken == "" {
		var err error
		*xSessionToken, err = getSessToken(token, authInf, xSessionToken)
		if err != nil {
			log.Fatal(err)
		}
		token = *xSessionToken
	}

	return token, nil

}

func newPostRequest() *postRequest {
	return &postRequest{
		method: "POST",
		url:    mafreebox + "api/v4/rrd/",
		header: "X-Fbx-App-Auth",
	}
}

// getDsl get dsl statistics
func getDsl(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int, error) {
	d := &database{
		DB:        "dsl",
		Fields:    []string{"rate_up", "rate_down", "snr_up", "snr_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int{}, err
	}
	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int{}, err
	}
	if resp.StatusCode == 404 {
		return []int{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		return []int{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int{}, err
		}
	}

	if rrdTest.ErrorCode != "" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getDsl")
		}
		return []int{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int{}, nil
	}

	result := []int{rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["snr_up"], rrdTest.Result.Data[0]["snr_down"]}
	return result, nil
}

// getTemp get temp statistics
func getTemp(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int, error) {
	d := &database{
		DB:        "temp",
		Fields:    []string{"cpum", "cpub", "sw", "hdd", "fan_speed"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, fmt.Sprintf(pr.url), buf)
	if err != nil {
		return []int{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int{}, err
	}
	if resp.StatusCode == 404 {
		return []int{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		return []int{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int{}, err
		}
	}

	if rrdTest.ErrorCode != "" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getTemp")
		}
		return []int{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int{}, nil
	}

	return []int{rrdTest.Result.Data[0]["cpum"], rrdTest.Result.Data[0]["cpub"], rrdTest.Result.Data[0]["sw"], rrdTest.Result.Data[0]["hdd"], rrdTest.Result.Data[0]["fan_speed"]}, nil
}

// getNet get net statistics
func getNet(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int, error) {
	d := &database{
		DB:        "net",
		Fields:    []string{"bw_up", "bw_down", "rate_up", "rate_down", "vpn_rate_up", "vpn_rate_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int{}, err
	}
	if resp.StatusCode == 404 {
		return []int{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		return []int{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int{}, err
		}
	}

	if rrdTest.ErrorCode != "" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getNet")
		}
		return []int{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int{}, nil
	}

	return []int{rrdTest.Result.Data[0]["bw_up"], rrdTest.Result.Data[0]["bw_down"], rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["vpn_rate_up"], rrdTest.Result.Data[0]["vpn_rate_down"]}, nil
}

// getSwitch get switch statistics
func getSwitch(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int, error) {
	d := &database{
		DB:        "switch",
		Fields:    []string{"rx_1", "tx_1", "rx_2", "tx_2", "rx_3", "tx_3", "rx_4", "tx_4"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int{}, err
	}
	if resp.StatusCode == 404 {
		return []int{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		return []int{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int{}, err
		}
	}

	if rrdTest.ErrorCode != "" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getSwitch")
		}
		return []int{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int{}, nil
	}

	return []int{rrdTest.Result.Data[0]["rx_1"], rrdTest.Result.Data[0]["tx_1"], rrdTest.Result.Data[0]["rx_2"], rrdTest.Result.Data[0]["tx_2"], rrdTest.Result.Data[0]["rx_3"], rrdTest.Result.Data[0]["tx_3"], rrdTest.Result.Data[0]["rx_4"], rrdTest.Result.Data[0]["tx_4"]}, nil
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
	if err != nil {
		return []lanHost{}, err
	}

	if lanResp.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []lanHost{}, err
		}
	}

	if lanResp.ErrorCode != "" {
		return []lanHost{}, lanResp.status()
	}

	return lanResp.Result, nil
}

// getLan get lan statistics
func getSystem(authInf *authInfo, pr *postRequest, xSessionToken *string) (system, error) {

	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return system{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return system{}, err
	}
	if resp.StatusCode == 404 {
		return system{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return system{}, err
	}

	systemResp := system{}
	err = json.Unmarshal(body, &systemResp)
	if err != nil {
		return system{}, err
	}

	return systemResp, nil
}
