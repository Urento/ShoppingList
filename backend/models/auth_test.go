package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestUpdateUser(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	newUsername := util.StringWithCharset(20)
	auth := Auth{
		EMail:    user.EMail,
		Username: newUsername,
	}

	err = auth.UpdateUser(user.EMail)
	if err != nil {
		t.Errorf("Error while updating user: %s", err)
	}

	u, err := GetUser(user.EMail)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	Equal(t, user.EMail, u.EMail)
	Equal(t, newUsername, u.Username)
}

func TestDisableAccount(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	err = DisableAccount(user.EMail)
	if err != nil {
		t.Errorf("Error while disabling account: %s", err)
	}

	disabled, err := IsDisabled(user.EMail)
	if err != nil {
		t.Errorf("Error while checking if account is disabled: %s", err)
	}

	Equal(t, disabled, true)
	Equal(t, nil, err)
}

func TestActivateAccount(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	err = DisableAccount(user.EMail)
	if err != nil {
		t.Errorf("Error while disabling account: %s", err)
	}

	disabledBefore, err := IsDisabled(user.EMail)
	if err != nil {
		t.Errorf("Error while checking if account is disabled: %s", err)
	}

	err = ActivateAccount(user.EMail)
	if err != nil {
		t.Errorf("Error while activating account: %s", err)
	}

	disabledAfter, err := IsDisabled(user.EMail)
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
		user, err := CreateUser()
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		u, err := GetUser(user.EMail)
		if err != nil {
			t.Errorf("Error while getting user: %s", err)
		}

		exists, err := ExistsUserID(u.ID)
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

func TestGetUserIDByEmail(t *testing.T) {
	Setup()

	user, err := CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	userId, err := GetUserIDByEmail(user.EMail)
	if err != nil {
		t.Errorf("Error while getting userid by email: %s", err)
	}

	Equal(t, user.ID, userId)
}

func TestResetPasswordFromUser(t *testing.T) {
	Setup()

	t.Run("Reset Password", func(t *testing.T) {
		username := util.StringWithCharset(5000)
		email := util.RandomEmail()
		password := util.StringWithCharset(5000)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, password, ip)
		if err != nil {
			t.Errorf("Error while creating account: %s", err)
		}

		newPassword := util.StringWithCharset(50000)
		err = ResetPasswordFromUser(email, newPassword, password, true)
		if err != nil {
			t.Errorf("Error while resetting password from user: %s", err)
		}

		Nil(t, err)
	})

	t.Run("Reset Password when the old password is wrong", func(t *testing.T) {
		username := util.StringWithCharset(5000)
		email := util.RandomEmail()
		password := util.StringWithCharset(5000)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, password, ip)
		if err != nil {
			t.Errorf("Error while creating account: %s", err)
		}

		newPassword := util.StringWithCharset(50000)
		err = ResetPasswordFromUser(email, newPassword, "dfbgjhdfbgdjfhbgdfg", true)

		Equal(t, "password is not correct", err.Error())
	})

	t.Run("Reset Password without old password", func(t *testing.T) {
		username := util.StringWithCharset(5000)
		email := util.RandomEmail()
		password := util.StringWithCharset(5000)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, password, ip)
		if err != nil {
			t.Errorf("Error while creating account: %s", err)
		}

		newPassword := util.StringWithCharset(50000)
		err = ResetPasswordFromUser(email, newPassword, "", false)
		if err != nil {
			t.Errorf("Error while resetting password from user: %s", err)
		}

		Nil(t, err)
	})
}
