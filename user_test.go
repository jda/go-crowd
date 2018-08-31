package crowd

import (
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {
	tv := PrepVars(t)
	c, err := New(tv.AppUsername, tv.AppPassword, tv.AppURL)
	if err != nil {
		t.Error(err)
	}

	user := os.Getenv("APP_USER_USERNAME")
	if user == "" {
		t.Skip("Can't run test because APP_USER_USERNAME undefined")
	}

	// test new session
	u, err := c.GetUser(user)
	if err != nil {
		t.Errorf("Error getting user info: %s\n", err)
	} else {
		t.Logf("Got user info: %+v\n", u)
	}

	if u.UserName == "" {
		t.Errorf("username was empty so we didn't get/decode a response from GetUser")
	}

}

func TestGetUserByEmail(t *testing.T) {
	tv := PrepVars(t)
	c, err := New(tv.AppUsername, tv.AppPassword, tv.AppURL)
	if err != nil {
		t.Error(err)
	}

	email := os.Getenv("APP_USER_EMAIL")
	if email == "" {
		t.Skip("Can't run test because APP_USER_EMAIL undefined")
	}

	// test get by email
	u, err := c.GetUserByEmail(email)
	if err != nil {
		t.Errorf("Error getting user info by email: %s\n", err)
	} else {
		t.Logf("Got user info by email: %+v\n", u)
	}

	if u.UserName == "" {
		t.Errorf("username was empty so we didn't get/decode a response from GetUserByEmail")
	}
}
