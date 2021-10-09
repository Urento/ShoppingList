package models

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestUpdateUser(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(20)
	ip := RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	newUsername := util.StringWithCharset(20)
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

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(20)
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

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(20)
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

func TestExistsUserId(t *testing.T) {
	Setup()

	t.Run("Exists User Id", func(t *testing.T) {
		pwd := util.StringWithCharset(20)
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(20)
		ip := RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		user, err := GetUser(email)
		if err != nil {
			t.Errorf("Error while getting user: %s", err)
		}

		exists, err := ExistsUserID(user.ID)
		if err != nil {
			t.Errorf("Error while checking if the userid exists: %s", err)
		}

		Equal(t, true, exists)
	})

	t.Run("Exists User Id when the user doesn't exist", func(t *testing.T) {
		exists, _ := ExistsUserID(999999999999999999)
		Equal(t, false, exists)
	})
}

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
