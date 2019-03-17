package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// getTemp get temp statistics
func getTemp() (int, int, int, int, int) {
	freeboxToken := os.Getenv("FREEBOX_TOKEN")
	if freeboxToken == "" {
		sessToken = getToken()
	}
	if sessToken == "" {
		sessToken = getSessToken(freeboxToken)
	}
	xtemp := database{
		DB:        "temp",
		Fields:    []string{"cpum", "cpub", "sw", "hdd", "fan_speed"},
		Precision: 10,
		DateStart: int(time.Now().Unix() - 10),
	}
	client := http.Client{}
	r, err := json.Marshal(xtemp)
	if err != nil {
		log.Fatalln(err)
	}
	buf := bytes.NewReader(r)
	req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/%s/rrd/", mafreebox, version), buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-Fbx-App-Auth", sessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 404 {
		log.Fatal(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &rrdTest)
	switch rrdTest.ErrorCode {
	case "auth_required":
		sessToken = getSessToken(freeboxToken)
	case "invalid_token":
		log.Fatalln("The app token you are trying to use is invalid or has been revoked")
	case "pending_token":
		log.Println("The app token you are trying to use has not been validated by user yet")
	case "insufficient_rights":
		log.Fatalln("Your app permissions does not allow accessing this API")
	case "denied_from_external_ip":
		log.Fatalln("You are trying to get an app_token from a remote IP")
	case "invalid_request":
		log.Fatalln("Your request is invalid")
	case "ratelimited":
		log.Fatalln("Too many auth error have been made from your IP")
	case "new_apps_denied":
		log.Fatalln("New application token request has been disabled")
	case "apps_denied":
		log.Fatalln("API access from apps has been disabled")
	case "internal_error":
		log.Fatalln("Internal error")
	}
	return rrdTest.Result.Data[0]["cpum"], rrdTest.Result.Data[0]["cpub"], rrdTest.Result.Data[0]["sw"], rrdTest.Result.Data[0]["hdd"], rrdTest.Result.Data[0]["fan_speed"]
}
