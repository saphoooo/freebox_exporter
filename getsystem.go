package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// getLan get lan statistics
func getSystem(authInf *authInfo) systemR {
	freeboxToken, err := setFreeboxToken(authInf)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%sapi/%s/system/", mafreebox, version), nil)
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
	err = json.Unmarshal(body, &systemResp)
	switch systemResp.ErrorCode {
	case "auth_required":
		sessToken, err = getSessToken(freeboxToken, authInf)
		if err != nil {
			log.Fatal(err)
		}
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
	return systemResp.Result
}
