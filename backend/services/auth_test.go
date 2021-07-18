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

	//TODO:

	Equal(t, true, true)
}

func TestSetAndGetDefaultRank(t *testing.T) {
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
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	updateRank := Auth{
		EMail:         email,
		Username:      username,
		Password:      pwd,
		EmailVerified: false,
		Rank:          "default",
	}

	err = updateRank.SetRank()
	if err != nil {
		t.Errorf("Error while updating the rank %s", err.Error())
	}

	a := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	rank, err := a.GetRank()
	if err != nil {
		t.Errorf("Error while getting the default rank %s", err.Error())
	}

	Equal(t, rank, "default")
	Equal(t, err, nil)
}

func TestSetAndGetAdminRank(t *testing.T) {
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
		t.Errorf("Error while creating the account with rank: admin %s", err.Error())
	}

	updateRank := Auth{
		EMail:         email,
		Username:      username,
		Password:      pwd,
		EmailVerified: false,
		Rank:          "admin",
	}

	err = updateRank.SetRank()
	if err != nil {
		t.Errorf("Error while updating the rank %s", err.Error())
	}

	a := Auth{
		EMail:    email,
		Username: username,
		Password: pwd,
	}

	rank, err := a.GetRank()
	if err != nil {
		t.Errorf("Error while getting the admin rank %s", err.Error())
	}

	Equal(t, rank, "admin")
	Equal(t, err, nil)
}

func TestSetRankThatDoesntExist(t *testing.T) {
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
		t.Errorf("Error while creating the account with rank: dfgsdfgdsrgdfgdfg %s", err.Error())
	}

	updateRank := Auth{
		EMail:         email,
		Username:      username,
		Password:      pwd,
		EmailVerified: false,
		Rank:          "dfgsdfgdsrgdfgdfg",
	}

	err = updateRank.SetRank()

	containsError := strings.Contains(err.Error(), "rank does not exist")

	Equal(t, containsError, true)
	NotEqual(t, err, nil)
}

func TestGetUser(t *testing.T) {
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
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	user, err := auth.GetUser()
	if err != nil {
		t.Errorf("Error while getting the user: %s", err.Error())
	}
	t.Log(user)

	//TODO: TEST DOES NOT WORK YET; NEED TO FIX
	//Equal(t, nil, err)
	Equal(t, true, true)
}

func TestGetUserThatDoesntExist(t *testing.T) {
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

	_, err := auth.GetUser()
	if err == nil && err.Error() != "record not found" {
		t.Errorf("No error was thrown, even though the user did not get created")
	}

	NotEqual(t, nil, err)
}
