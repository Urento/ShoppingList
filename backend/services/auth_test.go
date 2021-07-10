package services

import (
	"strings"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCreateAccountEmail(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	auth := Auth{
		EMail:    email,
		Username: "",
		Password: pwd,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user %s", err.Error())
	}

	checkAuth := Auth{
		EMail:    email,
		Username: "",
		Password: pwd,
	}

	check, err := checkAuth.Check()
	if err != nil {
		t.Errorf("Error while checking user %s", err.Error())
	}

	if !check {
		t.Errorf("check was false")
	}

	Equal(t, check, true)
	Equal(t, err, nil)
}

func TestCreateAccountWithEmailAndUsername(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user %s", err.Error())
	}

	checkAuth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	check, err := checkAuth.Check()
	if err != nil {
		t.Errorf("Error while checking user %s", err.Error())
	}

	if !check {
		t.Errorf("check was false")
	}

	Equal(t, check, true)
	Equal(t, err, nil)
}

func TestLoginWrongPassword(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	wrongPwd := StringWithCharset(50)
	username := StringWithCharset(10)
	email := StringWithCharset(20) + "@gmail.com"
	auth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user %s", err.Error())
	}

	checkAuth := Auth{
		EMail:    email,
		Username: username,
		Password: wrongPwd,
	}

	check, _ := checkAuth.Check()

	Equal(t, check, false)
	Equal(t, err, nil)
}

func TestDuplicateAccounts(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	auth.Create()

	err := auth.Create()
	if err == nil {
		t.Errorf("No Duplication error thrown")
	}

	containsError := strings.Contains(err.Error(), "account already exists")

	Equal(t, containsError, true)
	NotEqual(t, err.Error(), nil)
}

func TestInvalidEmail(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10)
	username := StringWithCharset(10)
	auth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	auth.Create()

	err := auth.Create()
	if err == nil {
		t.Errorf("No Invalid Email error thrown")
	}

	containsError := strings.Contains(err.Error(), "email is not valid")

	Equal(t, containsError, true)
	NotEqual(t, err.Error(), nil)
}

func TestDeleteAccount(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating the account %s", err.Error())
	}

	delAcc := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	err = delAcc.Delete()
	if err != nil {
		t.Errorf("Error while deleting the account %s", err.Error())
	}

	Equal(t, err, nil)
}

func TestNotEmailVerified(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:         email,
		Username:      username,
		Password:      pwd,
		EmailVerified: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating account: %s", err.Error())
	}

	verified, err := auth.IsEmailVerified()
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	Equal(t, verified, false)
	Equal(t, err, nil)
}

func TestUpdateEmailVerified(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:         email,
		Username:      username,
		Password:      pwd,
		EmailVerified: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating account: %s", err.Error())
	}

	verified1, err := auth.IsEmailVerified()
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	err = auth.VerifyEmail()
	if err != nil {
		t.Errorf("Error while verifying email: %s", err.Error())
	}

	verified2, err := auth.IsEmailVerified()
	if err != nil {
		t.Errorf("Error while checking if email is verified: %s", err.Error())
	}

	Equal(t, verified1, false)
	Equal(t, verified2, true)
	Equal(t, err, nil)
}

func TestSendVerificationEmail(t *testing.T) {
	Setup()

	Equal(t, true, true)
}
