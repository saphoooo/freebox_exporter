package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// storeToken stores app_token in ~/.freebox_token
func storeToken(token string, authInf *authInfo) error {
	if token == "" {
		return errors.New("token should not be blank")
	}

	err := os.Setenv("FREEBOX_TOKEN", token)
	if err != nil {
		return err
	}

	if _, err := os.Stat(authInf.myStore.location); os.IsNotExist(err) {
		err := ioutil.WriteFile(authInf.myStore.location, []byte(token), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

// retreiveToken gets the token from file and
// load it in environment variable
func retreiveToken(authInf *authInfo) (string, error) {
	if _, err := os.Stat(authInf.myStore.location); os.IsNotExist(err) {
		return "", err
	}
	data, err := ioutil.ReadFile(authInf.myStore.location)
	if err != nil {
		return "", err
	}
	err = os.Setenv("FREEBOX_TOKEN", string(data))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getTrackID is the initial request to freebox API
// get app_token and track_id
func getTrackID(authInf *authInfo) (*track, error) {

	req, _ := json.Marshal(authInf.myApp)
	buf := bytes.NewReader(req)
	//resp, err := http.Post(mafreebox+"api/"+version+"/login/authorize/", "application/json", buf)
	resp, err := http.Post(authInf.myAPI.authz, "application/json", buf)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	trackID := track{}
	err = json.Unmarshal(body, &trackID)
	if err != nil {
		return nil, err
	}

	err = storeToken(trackID.Result.AppToken, authInf)
	if err != nil {
		return nil, err
	}

	return &trackID, nil
}

// getGranted waits for user to validate on the freebox front panel
// with a timeout of 15 seconds
func getGranted(authInf *authInfo) error {
	trackID, err := getTrackID(authInf)
	if err != nil {
		return err
	}
	//url := mafreebox + "api/" + version + "/login/authorize/" + strconv.Itoa(trackID.Result.TrackID)
	url := authInf.myAPI.authz + strconv.Itoa(trackID.Result.TrackID)
	for i := 0; i < 15; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &granted)
		if err != nil {
			return err
		}

		switch granted.Result.Status {
		case "unknown":
			return errors.New("the app_token is invalid or has been revoked")
		case "pending":
			log.Println("the user has not confirmed the authorization request yet")
		case "timeout":
			return errors.New("the user did not confirmed the authorization within the given time")
		case "granted":
			log.Println("the app_token is valid and can be used to open a session")
			i = 15
		case "denied":
			return errors.New("the user denied the authorization request")
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// getChallenge makes sure the app always has a valid challenge
func getChallenge(authInf *authInfo) error {
	//resp, err := http.Get(mafreebox + "api/" + version + "/login/")
	resp, err := http.Get(authInf.myAPI.login)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &challenged)
	if err != nil {
		return err
	}
	return nil
}

// hmacSha1 encodes app_token in hmac-sha1 and stores it in password
func hmacSha1(appToken, challenge string) string {
	//secret := []byte(appToken)
	//message := []byte(challenge)
	//hash := hmac.New(sha1.New, secret)
	//hash.Write(message)
	hash := hmac.New(sha1.New, []byte(appToken))
	hash.Write([]byte(challenge))
	return hex.EncodeToString(hash.Sum(nil))
}

// getSession gets a session with freeebox API
func getSession(authInf *authInfo, passwd string) (*sessionToken, error) {
	s := session{
		AppID:    authInf.myApp.AppID,
		Password: passwd,
	}
	req, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(req)
	//resp, err := http.Post(mafreebox+"api/"+version+"/login/session/", "application/json", buf)
	resp, err := http.Post(authInf.myAPI.loginSession, "application/json", buf)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := sessionToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// getToken gets a valid session_token and asks for user to change
// the set of permissions on the API
func getToken(authInf *authInfo) (string, error) {
	if _, err := os.Stat(authInf.myStore.location); os.IsNotExist(err) {
		err = getGranted(authInf)
		if err != nil {
			return "", err
		}
		reader := bufio.NewReader(os.Stdin)
		log.Println("check \"Modification des réglages de la Freebox\" and press enter")
		_, _ = reader.ReadString('\n')
	} else {
		_, err := retreiveToken(authInf)
		if err != nil {
			return "", err
		}
	}
	err := getChallenge(authInf)
	if err != nil {
		return "", err
	}
	password := hmacSha1(os.Getenv("FREEBOX_TOKEN"), challenged.Result.Challenge)
	token, err := getSession(authInf, password)
	if err != nil {
		return "", err
	}
	if token.Success {
		fmt.Println("successfully authenticated")
	} else {
		return "", errors.New(token.Msg)
	}
	return token.Result.SessionToken, nil
}

// getSessToken gets a new token session when the old one has expired
func getSessToken(t string, authInf *authInfo) (string, error) {
	err := getChallenge(authInf)
	if err != nil {
		return "", err
	}
	password := hmacSha1(t, challenged.Result.Challenge)
	token, err := getSession(authInf, password)
	if err != nil {
		return "", err
	}
	if token.Success == false {
		return "", errors.New(token.Msg)
	}
	return token.Result.SessionToken, nil
}
