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
		"db_error":                errors.New("Oops, the database you are trying to access doesn't seem to exist"),
		"nodev":                   errors.New("Invalid interface"),
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

func getConnectionXdsl(authInf *authInfo, pr *postRequest, xSessionToken *string) (connectionXdsl, error) {
	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return connectionXdsl{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return connectionXdsl{}, err
	}
	if resp.StatusCode == 404 {
		return connectionXdsl{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return connectionXdsl{}, err
	}

	connectionXdslResp := connectionXdsl{}
	err = json.Unmarshal(body, &connectionXdslResp)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return connectionXdsl{}, err
	}

	return connectionXdslResp, nil
}

// getDsl get dsl statistics
func getDsl(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	d := &database{
		DB:        "dsl",
		Fields:    []string{"rate_up", "rate_down", "snr_up", "snr_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int64{}, err
	}
	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int64{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int64{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int64{}, err
	}
	if resp.StatusCode == 404 {
		return []int64{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int64{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return []int64{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int64{}, err
		}
	}

	if rrdTest.ErrorCode != "" && rrdTest.ErrorCode != "auth_required" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getDsl")
		}
		return []int64{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int64{}, nil
	}

	result := []int64{rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["snr_up"], rrdTest.Result.Data[0]["snr_down"]}
	return result, nil
}

// getTemp get temp statistics
func getTemp(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	d := &database{
		DB:        "temp",
		Fields:    []string{"cpum", "cpub", "sw", "hdd", "fan_speed"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int64{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int64{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, fmt.Sprintf(pr.url), buf)
	if err != nil {
		return []int64{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int64{}, err
	}
	if resp.StatusCode == 404 {
		return []int64{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int64{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return []int64{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int64{}, err
		}
	}

	if rrdTest.ErrorCode != "" && rrdTest.ErrorCode != "auth_required" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getTemp")
		}
		return []int64{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int64{}, nil
	}

	return []int64{rrdTest.Result.Data[0]["cpum"], rrdTest.Result.Data[0]["cpub"], rrdTest.Result.Data[0]["sw"], rrdTest.Result.Data[0]["hdd"], rrdTest.Result.Data[0]["fan_speed"]}, nil
}

// getNet get net statistics
func getNet(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	d := &database{
		DB:        "net",
		Fields:    []string{"bw_up", "bw_down", "rate_up", "rate_down", "vpn_rate_up", "vpn_rate_down"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int64{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int64{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int64{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int64{}, err
	}
	if resp.StatusCode == 404 {
		return []int64{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int64{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return []int64{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int64{}, err
		}
	}

	if rrdTest.ErrorCode != "" && rrdTest.ErrorCode != "auth_required" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getNet")
		}
		return []int64{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int64{}, nil
	}

	return []int64{rrdTest.Result.Data[0]["bw_up"], rrdTest.Result.Data[0]["bw_down"], rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["vpn_rate_up"], rrdTest.Result.Data[0]["vpn_rate_down"]}, nil
}

// getSwitch get switch statistics
func getSwitch(authInf *authInfo, pr *postRequest, xSessionToken *string) ([]int64, error) {
	d := &database{
		DB:        "switch",
		Fields:    []string{"rx_1", "tx_1", "rx_2", "tx_2", "rx_3", "tx_3", "rx_4", "tx_4"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}

	freeboxToken, err := setFreeboxToken(authInf, xSessionToken)
	if err != nil {
		return []int64{}, err
	}

	client := http.Client{}
	r, err := json.Marshal(*d)
	if err != nil {
		return []int64{}, err
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		return []int64{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return []int64{}, err
	}
	if resp.StatusCode == 404 {
		return []int64{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []int64{}, err
	}
	rrdTest := rrd{}
	err = json.Unmarshal(body, &rrdTest)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return []int64{}, err
	}

	if rrdTest.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []int64{}, err
		}
	}

	if rrdTest.ErrorCode != "" && rrdTest.ErrorCode != "auth_required" {
		if rrdTest.status().Error() == "Unknown return code from the API" {
			fmt.Println("getSwitch")
		}
		return []int64{}, rrdTest.status()
	}

	if len(rrdTest.Result.Data) == 0 {
		return []int64{}, nil
	}

	return []int64{rrdTest.Result.Data[0]["rx_1"], rrdTest.Result.Data[0]["tx_1"], rrdTest.Result.Data[0]["rx_2"], rrdTest.Result.Data[0]["tx_2"], rrdTest.Result.Data[0]["rx_3"], rrdTest.Result.Data[0]["tx_3"], rrdTest.Result.Data[0]["rx_4"], rrdTest.Result.Data[0]["tx_4"]}, nil
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
		if debug {
			log.Println(string(body))
		}
		return []lanHost{}, err
	}

	if lanResp.ErrorCode == "auth_required" {
		*xSessionToken, err = getSessToken(freeboxToken, authInf, xSessionToken)
		if err != nil {
			return []lanHost{}, err
		}
	}

	if lanResp.ErrorCode != "" && lanResp.ErrorCode != "auth_required" {
		return []lanHost{}, lanResp.status()
	}

	return lanResp.Result, nil
}

func getFreeplug(authInf *authInfo, pr *postRequest, xSessionToken *string) (freeplug, error) {
	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return freeplug{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return freeplug{}, err
	}
	if resp.StatusCode == 404 {
		return freeplug{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return freeplug{}, err
	}

	freeplugResp := freeplug{}
	err = json.Unmarshal(body, &freeplugResp)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return freeplug{}, err
	}

	return freeplugResp, nil
}

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
		if debug {
			log.Println(string(body))
		}
		return system{}, err
	}

	return systemResp, nil
}

func getWifi(authInf *authInfo, pr *postRequest, xSessionToken *string) (wifi, error) {
	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return wifi{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return wifi{}, err
	}
	if resp.StatusCode == 404 {
		return wifi{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wifi{}, err
	}

	wifiResp := wifi{}
	err = json.Unmarshal(body, &wifiResp)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return wifi{}, err
	}

	return wifiResp, nil
}

func getWifiStations(authInf *authInfo, pr *postRequest, xSessionToken *string) (wifiStations, error) {
	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return wifiStations{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return wifiStations{}, err
	}
	if resp.StatusCode == 404 {
		return wifiStations{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return wifiStations{}, err
	}

	wifiStationResp := wifiStations{}
	err = json.Unmarshal(body, &wifiStationResp)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return wifiStations{}, err
	}

	return wifiStationResp, nil
}

func getVpnServer(authInf *authInfo, pr *postRequest, xSessionToken *string) (vpnServer, error) {
	client := http.Client{}
	req, err := http.NewRequest(pr.method, pr.url, nil)
	if err != nil {
		return vpnServer{}, err
	}
	req.Header.Add(pr.header, *xSessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return vpnServer{}, err
	}
	if resp.StatusCode == 404 {
		return vpnServer{}, errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vpnServer{}, err
	}

	vpnServerResp := vpnServer{}
	err = json.Unmarshal(body, &vpnServerResp)
	if err != nil {
		if debug {
			log.Println(string(body))
		}
		return vpnServer{}, err
	}

	return vpnServerResp, nil
}
