package models

import (
	"strings"
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

func TestCreateAccountEmail(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	ip := util.RandomIPAddress()

	t.Run("Create Account Without IP", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(20)

		err := CreateAccount(email, username, pwd, "")
		if err != nil {
			t.Errorf("Error while creating user %s", err.Error())
		}

		check, err := CheckAuth(email, pwd, "")
		if err != nil {
			t.Errorf("Error while checking user %s", err.Error())
		}

		if !check {
			t.Errorf("check was false")
		}

		Equal(t, true, check)
		Equal(t, nil, err)
	})

	t.Run("Create Account With IP", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(20)

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating user %s", err.Error())
		}

		ip2 := util.RandomIPAddress()

		check, err := CheckAuth(email, pwd, ip2)
		if err != nil {
			t.Errorf("Error while checking user %s", err.Error())
		}

		if !check {
			t.Errorf("check was false")
		}

		Equal(t, true, check)
		Equal(t, nil, err)
	})

	t.Run("Create Account when the Email already exist", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(10)

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating account: %s", err)
		}

		err = CreateAccount(email, username, pwd, ip)
		if err == nil {
			t.Errorf("No Duplication error thrown")
		}

		containsError := strings.Contains(err.Error(), "email is already being used")

		Equal(t, true, containsError)
		NotEqual(t, nil, err)
	})

	t.Run("Create Account with invalid email", func(t *testing.T) {
		email := util.StringWithCharset(10)
		username := util.StringWithCharset(10)

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating account: %s", err)
		}

		err = CreateAccount(email, username, pwd, ip)
		if err == nil {
			t.Errorf("No Invalid Email error thrown")
		}

		containsError := strings.Contains(err.Error(), "email is not valid")

		Equal(t, true, containsError)
		NotEqual(t, nil, err)
	})
}

func TestCheckAuth(t *testing.T) {
	Setup()

	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(20)
	pwd := util.StringWithCharset(20)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating account: %s", err)
	}

	t.Run("When Account is Disabled", func(t *testing.T) {
		err = DisableAccount(email)
		if err != nil {
			t.Errorf("Error while disabling account: %s", err)
		}

		disabledBefore, err := IsDisabled(email)
		if err != nil {
			t.Errorf("Error while checking if the account is disabled: %s", err)
		}

		ok, err := CheckAuth(email, pwd, ip)
		if err != nil {
			t.Errorf("Error while checking auth: %s", err)
		}

		disabledAfter, err := IsDisabled(email)
		if err != nil {
			t.Errorf("Error while checking if the account is disabled: %s", err)
		}

		err = ActivateAccount(email)
		if err != nil {
			t.Errorf("Error while activating account: %s", err)
		}

		Equal(t, true, disabledBefore)
		Equal(t, true, disabledAfter)
		Equal(t, true, ok)
		Equal(t, nil, err)
	})

	t.Run("Normal Login", func(t *testing.T) {
		ok, err := CheckAuth(email, pwd, ip)
		if err != nil {
			t.Errorf("Error while checking auth: %s", err)
		}

		disabled, err := IsDisabled(email)
		if err != nil {
			t.Errorf("Error while checking if the account is disabled: %s", err)
		}

		Equal(t, false, disabled)
		Equal(t, true, ok)
		Equal(t, nil, err)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		ok, _ := CheckAuth(email, "dsjkhfgjkhdsbvsdfg", "")

		Equal(t, false, ok)
	})

	t.Run("Wrong Email and Password", func(t *testing.T) {
		ok, _ := CheckAuth("kjfdghjdfbgjhdfbg@gmail.com", "khdfbgkjhdfgbhjdfgbdf", "")

		Equal(t, false, ok)
	})
}

func TestDeleteAccount(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(10)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating the account %s", err.Error())
	}

	err = DeleteAccount(email, pwd)
	if err != nil {
		t.Errorf("Error while deleting the account %s", err.Error())
	}

	Equal(t, nil, err)
}

func TestEmailVerified(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(10)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating account: %s", err.Error())
	}

	verified, err := IsEmailVerified(email)
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	Equal(t, false, verified)
	Equal(t, nil, err)
}

func TestUpdateEmailVerified(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(10)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating account: %s", err.Error())
	}

	verified1, err := IsEmailVerified(email)
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	err = VerifyEmail(email)
	if err != nil {
		t.Errorf("Error while verifying email: %s", err.Error())
	}

	verified2, err := IsEmailVerified(email)
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	Equal(t, nil, err)
	Equal(t, false, verified1)
	Equal(t, true, verified2)
}

func TestSetAndGetRank(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)

	t.Run("Default Rank", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(10)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating the account with rank: default %s", err.Error())
		}

		err = SetRank(email, "default")
		if err != nil {
			t.Errorf("Error while updating the rank %s", err.Error())
		}

		rank, err := GetRank(email)
		if err != nil {
			t.Errorf("Error while getting the default rank %s", err.Error())
		}

		Equal(t, "default", rank)
		Equal(t, nil, err)
	})

	t.Run("Admin Rank", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(10)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating the account with rank: admin %s", err.Error())
		}

		err = SetRank(email, "admin")
		if err != nil {
			t.Errorf("Error while updating the rank %s", err.Error())
		}

		rank, err := GetRank(email)
		if err != nil {
			t.Errorf("Error while getting the admin rank %s", err.Error())
		}

		Equal(t, "admin", rank)
		Equal(t, nil, err)
	})

	t.Run("Rank does not exist", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(10)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating the account with rank: dfgsdfgdsrgdfgdfg %s", err.Error())
		}

		err = SetRank(email, "dkhfgbjhdfbgjhdfg")

		containsError := strings.Contains(err.Error(), "rank does not exist")

		Equal(t, true, containsError)
		NotEqual(t, nil, err)
	})
}

func TestGetUser(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)

	t.Run("Get User", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(10)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating the account with rank: default %s", err.Error())
		}

		user, err := GetUser(email)
		if err != nil {
			t.Errorf("Error while getting the user: %s", err.Error())
		}

		if user.ID <= 0 {
			t.Errorf("ID is 0, expected is a number higher or equal to 1")
		}

		Equal(t, nil, err)
		NotEqual(t, nil, user)
		Equal(t, email, user.EMail)
		Equal(t, username, user.Username)
		Equal(t, false, user.EmailVerified)
		Equal(t, "default", user.Rank)
	})

	t.Run("Get User that doesn't exist", func(t *testing.T) {
		email := util.StringWithCharset(100) + "@gmail.com"

		_, err := GetUser(email)
		if err == nil && err.Error() != "record not found" {
			t.Errorf("No error was thrown, even though the user did not get created")
		}

		NotEqual(t, nil, err)
	})
}

func TestEnableTwoFactorAuthentication(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(10)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	err = SetTwoFactorAuthentication(email, true)
	if err != nil {
		t.Errorf("Error while updating Two Factor Authentication Status %s", err)
	}

	isEnabled, err := IsTwoFactorEnabled(email)
	if err != nil {
		t.Errorf("Error while getting Two Factor Authentication Status %s", err)
	}

	Equal(t, nil, err)
	Equal(t, true, isEnabled)
}

func TestUpdateIPAndGet(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)
	email := util.StringWithCharset(10) + "@gmail.com"
	username := util.StringWithCharset(100)
	ip := util.RandomIPAddress()

	err := CreateAccount(email, username, pwd, ip)
	if err != nil {
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	ip2 := util.RandomIPAddress()

	err = UpdateIP(email, ip2)
	if err != nil {
		t.Errorf("Error while updating ip: %s", err)
	}

	newIP, err := GetIP(email)
	if err != nil {
		t.Errorf("Error while getting ip: %s", err)
	}

	Equal(t, ip, newIP)
	Equal(t, nil, err)
}

func TestUpdateUsername(t *testing.T) {
	Setup()

	pwd := util.StringWithCharset(20)

	t.Run("Update Username", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(100)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		username2 := util.StringWithCharset(30)

		err = SetUsername(email, username2)
		if err != nil {
			t.Errorf("Error while updating username: %s", err)
		}

		uName, err := GetUsername(email)
		if err != nil {
			t.Errorf("Error while getting username: %s", err)
		}

		Equal(t, username, uName)
		Equal(t, nil, err)
	})

	t.Run("Update Username with over 32 characters", func(t *testing.T) {
		email := util.StringWithCharset(10) + "@gmail.com"
		username := util.StringWithCharset(100)
		ip := util.RandomIPAddress()

		err := CreateAccount(email, username, pwd, ip)
		if err != nil {
			t.Errorf("Error while creating user: %s", err)
		}

		username2 := util.StringWithCharset(34)

		err = SetUsername(email, username2)
		if err == nil {
			t.Errorf("No Error thrown even though the username is over 32 characters")
		}

		Equal(t, "username can only be a maximum of 32 characters long", err.Error())
		NotEqual(t, nil, err)
	})
}
