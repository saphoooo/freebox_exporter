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
