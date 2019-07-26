package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// storeToken stores app_token in ~/.freebox_token
func storeToken(t string) {
	err := os.Setenv("FREEBOX_TOKEN", t)
	if err != nil {
		log.Fatalln(err)
	}
	home := os.Getenv("HOME")
	fileStore := home + "/.freebox_token"
	if _, err := os.Stat(fileStore); os.IsNotExist(err) {
		err := ioutil.WriteFile(fileStore, []byte(t), 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println(fileStore, "already exists")
	}
}

// retreiveToken gets the token from file and
// load it in environment variable
func retreiveToken() {
	home := os.Getenv("HOME")
	fileStore := home + "/.freebox_token"
	if _, err := os.Stat(fileStore); os.IsNotExist(err) {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile(fileStore)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("FREEBOX_TOKEN", string(data))
	if err != nil {
		log.Fatal(err)
	}
}

// getTrackID is the initial request to freebox API
// get app_token and track_id
func getTrackID(app *app, fb *freebox) error {

	req, _ := json.Marshal(app)
	buf := bytes.NewReader(req)
	//resp, err := http.Post(mafreebox+"api/"+version+"/login/authorize/", "application/json", buf)
	resp, err := http.Post(fb.uri, "application/json", buf)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &trackID)
	if err != nil {
		return err
	}
	storeToken(trackID.Result.AppToken)
	return nil
}

// getGranted waits for user to validate on the freebox front panel
// with a timeout of 15 seconds
func getGranted(fb *freebox) {
	err := getTrackID(promExporter, fb)
	if err != nil {
		log.Fatalln(err)
	}
	url := mafreebox + "api/" + version + "/login/authorize/" + strconv.Itoa(trackID.Result.TrackID)
	for i := 0; i < 15; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &granted)
		if err != nil {
			log.Fatalln(err)
		}
		switch granted.Result.Status {
		case "unknown":
			log.Println("the app_token is invalid or has been revoked")
			os.Exit(1)
		case "pending":
			log.Println("the user has not confirmed the authorization request yet")
		case "timeout":
			log.Println("the user did not confirmed the authorization within the given time")
			os.Exit(1)
		case "granted":
			log.Println("the app_token is valid and can be used to open a session")
			i = 15
		case "denied":
			log.Println("the user denied the authorization request")
			os.Exit(1)
		}
		time.Sleep(1 * time.Second)
	}
}

// getChallenge makes sure the app always has a valid challenge
func getChallenge() {
	resp, err := http.Get(mafreebox + "api/" + version + "/login/")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &challenged)
	if err != nil {
		log.Fatalln(err)
	}
}

// hmacSha1 encodes app_token in hmac-sha1 and stores it in password
func hmacSha1(appToken, challenge string) string {
	secret := []byte(appToken)
	message := []byte(challenge)
	hash := hmac.New(sha1.New, secret)
	hash.Write(message)
	return hex.EncodeToString(hash.Sum(nil))
}

// getSession gets a session with freeebox API
func getSession(app, passwd string) {
	s := session{
		AppID:    app,
		Password: passwd,
	}
	req, _ := json.Marshal(s)
	buf := bytes.NewReader(req)
	resp, err := http.Post(mafreebox+"api/"+version+"/login/session/", "application/json", buf)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatalln(err)
	}
}

// getToken gets a valid session_token and asks for user to change
// the set of permissions on the API
func getToken() string {
	home := os.Getenv("HOME")
	fileStore := home + "/.freebox_token"
	if _, err := os.Stat(fileStore); os.IsNotExist(err) {
		fb := &freebox{
			uri: mafreebox + "api/" + version + "/login/authorize/",
		}
		getGranted(fb)
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("check \"Modification des réglages de la Freebox\" and press enter")
		_, _ = reader.ReadString('\n')
	} else {
		retreiveToken()
	}
	getChallenge()
	password := hmacSha1(os.Getenv("FREEBOX_TOKEN"), challenged.Result.Challenge)
	getSession(promExporter.AppID, password)
	if token.Success {
		fmt.Println("successfully authenticated")
	} else {
		log.Fatal(token.Msg)
	}
	return token.Result.SessionToken
}

// getSessToken gets a new token session when the old one has expired
func getSessToken(t string) string {
	getChallenge()
	password := hmacSha1(t, challenged.Result.Challenge)
	getSession(promExporter.AppID, password)
	if token.Success == false {
		log.Fatal(token.Msg)
	}
	return token.Result.SessionToken
}
