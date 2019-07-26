package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
func (d *database) getDsl(fb *freebox, st *store, pr *postRequest) (int, int, int, int) {
	freeboxToken := setFreeboxToken(fb, st)
	client := http.Client{}
	r, err := json.Marshal(d)
	if err != nil {
		log.Fatalln(err)
	}
	buf := bytes.NewReader(r)
	//req, err := http.NewRequest("POST", fmt.Sprintf("%sapi/%s/rrd/", mafreebox, version), buf)
	req, err := http.NewRequest(pr.method, pr.url, buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add(pr.header, sessToken)
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
	if len(rrdTest.Result.Data) == 0 {
		return 0, 0, 0, 0
	}
	return rrdTest.Result.Data[0]["rate_up"], rrdTest.Result.Data[0]["rate_down"], rrdTest.Result.Data[0]["snr_up"], rrdTest.Result.Data[0]["snr_down"]
}
