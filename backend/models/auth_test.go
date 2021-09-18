package models

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(20)
	ip := RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	newUsername := StringWithCharset(20)
	auth := Auth{
		EMail:    email,
		Username: newUsername,
	}

	err = auth.UpdateUser(email)
	if err != nil {
		t.Errorf("Error while updating user: %s", err)
	}

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	Equal(t, email, user.EMail)
	Equal(t, newUsername, user.Username)
}

func TestDisableAccount(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(20)
	ip := RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	err = DisableAccount(email)
	if err != nil {
		t.Errorf("Error while disabling account: %s", err)
	}

	disabled, err := IsDisabled(email)
	if err != nil {
		t.Errorf("Error while checking if account is disabled: %s", err)
	}

	Equal(t, disabled, true)
	Equal(t, nil, err)
}

func TestActivateAccount(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(20)
	ip := RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	err = DisableAccount(email)
	if err != nil {
		t.Errorf("Error while disabling account: %s", err)
	}

	disabledBefore, err := IsDisabled(email)
	if err != nil {
		t.Errorf("Error while checking if account is disabled: %s", err)
	}

	err = ActivateAccount(email)
	if err != nil {
		t.Errorf("Error while activating account: %s", err)
	}

	disabledAfter, err := IsDisabled(email)
	if err != nil {
		t.Errorf("Error while checking if account is disabled: %s", err)
	}

	Equal(t, disabledBefore, true)
	Equal(t, disabledAfter, false)
	Equal(t, nil, err)
}

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
