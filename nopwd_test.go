package nopwd

import (
	"strings"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	np := NewNoPwd("test.com", "abracadabra", 10)
	link, err := np.GenerateLoginLink("richard@place.com")
	if err != nil {
		println(err)
		t.Errorf("error generting code")
	}
	parts := strings.Split(link, "=")
	code := parts[1]
	valid, email, err := np.ValidateCode(code)
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

	// test fails after expiration
	valid, email, err = np.validateCodeAtTime(code, time.Now().Unix()+6000)
	if err == nil {
		t.Errorf("should have failed expiration")
	}
	if valid != false {
		t.Errorf("code should not be valid")
	}

	// test fails if bad code
	valid, email, err = np.validateCodeAtTime(code+"fail", time.Now().Unix())
	if err == nil {
		t.Errorf("should have failed expiration")
	}
	if valid != false {
		t.Errorf("code should not be valid")
	}

	// test succeeds a little bit of time before expiration
	valid, email, err = np.validateCodeAtTime(code, time.Now().Unix()+20)
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
}
