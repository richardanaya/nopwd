package nopwd

import (
	"strings"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	np := NewNoPwd("abracadabra")
	link, err := np.GenerateLoginLink("test.com/login_code=", "richard@place.com", 10)
	if err != nil {
		println(err)
		t.Errorf("error generating login code")
	}
	parts := strings.Split(link, "=")
	code := parts[1]
	valid, email, err := np.ValidateLoginCode(code)
	if err != nil {
		println(err.Error())
		t.Errorf("unexpected code validation failure")
		return
	}
	if email != "richard@place.com" {
		t.Errorf("wrong email")
	}
	if valid != true {
		t.Errorf("code not valid")
	}

	_, _, err = np.ValidateAPICode(code)
	if err == nil {
		t.Errorf("should not be able to validate login as a api code")
	}

	// test fails after expiration
	valid, email, err = np.validateCodeAtTime(code, "login", time.Now().Unix()+6000)
	if err == nil {
		t.Errorf("should have failed expiration")
	}
	if valid != false {
		t.Errorf("code should not be valid")
	}

	// test fails if bad code
	valid, email, err = np.validateCodeAtTime(code+"fail", "login", time.Now().Unix())
	if err == nil {
		t.Errorf("should have failed expiration")
	}
	if valid != false {
		t.Errorf("code should not be valid")
	}

	// test succeeds a little bit of time before expiration
	valid, email, err = np.validateCodeAtTime(code, "login", time.Now().Unix()+20)
	if email != "richard@place.com" {
		t.Errorf("wrong email")
	}
	if err != nil {
		println(err)
		t.Errorf("unexpected code validation failure")
	}
	if valid != true {
		t.Errorf("code not valid")
	}

	apiCode, err := np.GenerateAPICode("richard@place.com", 10)
	if err != nil {
		println(err)
		t.Errorf("error generating api code")
	}

	valid, email, err = np.ValidateAPICode(apiCode)
	if err != nil {
		println(err.Error())
		t.Errorf("unexpected code validation failure")
		return
	}
	if email != "richard@place.com" {
		t.Errorf("wrong email")
	}
	if valid != true {
		t.Errorf("code not valid")
	}
}
