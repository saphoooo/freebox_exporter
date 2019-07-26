package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRetreiveToken(t *testing.T) {
	st := &store{
		location: "/tmp/token",
	}

	_, err := retreiveToken(st)
	if err.Error() != "stat /tmp/token: no such file or directory" {
		t.Error("Expected bla, but got", err)
	}

	ioutil.WriteFile(st.location, []byte("IOI"), 0600)
	defer os.Remove(st.location)

	token, err := retreiveToken(st)
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
	st := &store{}

	err := storeToken(token, st)
	if err.Error() != "token should not be blank" {
		t.Error("Expected token should not be blank, but got", err)
	}

	token = "IOI"
	err = storeToken(token, st)
	if err.Error() != "open : no such file or directory" {
		t.Error("Expected open : no such file or directory, but got", err)
	}

	st.location = "/tmp/token"
	err = storeToken(token, st)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}
	defer os.Remove(st.location)

	token = os.Getenv("FREEBOX_TOKEN")
	if token != "IOI" {
		t.Error("Expected IOI, but got", token)
	}
	os.Unsetenv("FREEBOX_TOKEN")

	data, err := ioutil.ReadFile(st.location)
	if err != nil {
		t.Error("Expected no err, but got", err)
	}

	if string(data) != "IOI" {
		t.Error("Expected IOI, but got", string(data))
	}

}
