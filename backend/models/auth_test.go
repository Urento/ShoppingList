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
	t.Log(email)

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

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
