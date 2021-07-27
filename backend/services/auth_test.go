package services

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
)

func TestCreateAccountEmail(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(20)
	ip := RandomIPAddress()
	auth := Auth{
		EMail:     email,
		Username:  username,
		Password:  pwd,
		IPAddress: ip,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user %s", err.Error())
	}

	ip2 := RandomIPAddress()
	checkAuth := Auth{
		EMail:     email,
		Password:  pwd,
		IPAddress: ip2,
	}

	check, err := checkAuth.Check()
	if err != nil {
		t.Errorf("Error while checking user %s", err.Error())
	}

	if !check {
		t.Errorf("check was false")
	}

	Equal(t, true, check)
	Equal(t, nil, err)
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

	Equal(t, true, check)
	Equal(t, nil, err)
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

	Equal(t, false, check)
	Equal(t, nil, err)
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

	Equal(t, true, containsError)
	NotEqual(t, nil, err)
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

	Equal(t, true, containsError)
	NotEqual(t, nil, err)
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

	Equal(t, nil, err)
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

	Equal(t, false, verified)
	Equal(t, nil, err)
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

	Equal(t, nil, err)
	Equal(t, false, verified1)
	Equal(t, true, verified2)
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

	Equal(t, "default", rank)
	Equal(t, nil, err)
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

	Equal(t, "admin", rank)
	Equal(t, nil, err)
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

	Equal(t, true, containsError)
	NotEqual(t, nil, err)
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

	if user.ID <= 0 {
		t.Errorf("ID is 0, expected is a number higher or equal to 1")
	}

	Equal(t, nil, err)
	NotEqual(t, nil, user)
	Equal(t, email, user.EMail)
	Equal(t, username, user.Username)
	Equal(t, false, user.EmailVerified)
	Equal(t, "default", user.Rank)
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

func TestEnableTwoFactorAuthentication(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	auth := Auth{
		EMail:                   email,
		Username:                username,
		Password:                pwd,
		EmailVerified:           false,
		TwoFactorAuthentication: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	twoFAUser := Auth{
		EMail:                   email,
		Username:                username,
		Password:                pwd,
		EmailVerified:           false,
		TwoFactorAuthentication: true,
	}

	err = twoFAUser.SetTwoFactorAuthentication()
	if err != nil {
		t.Errorf("Error while updating Two Factor Authentication Status %s", err)
	}

	isEnabled, err := auth.IsTwoFactorEnabled()
	if err != nil {
		t.Errorf("Error while getting Two Factor Authentication Status %s", err)
	}

	Equal(t, nil, err)
	Equal(t, true, isEnabled)
}

func TestUpdateIPAndGet(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	auth := Auth{
		EMail:                   email,
		Password:                pwd,
		EmailVerified:           false,
		TwoFactorAuthentication: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating the account with rank: default %s", err.Error())
	}

	ip := RandomIPAddress()
	updateIpAuthObj := Auth{
		EMail:     email,
		IPAddress: ip,
	}

	err = updateIpAuthObj.UpdateIP()
	if err != nil {
		t.Errorf("Error while updating ip: %s", err)
	}

	newIP, err := auth.GetIP()
	if err != nil {
		t.Errorf("Error while getting ip: %s", err)
	}

	Equal(t, ip, newIP)
	Equal(t, nil, err)
}

func TestUpdateUsername(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	auth := Auth{
		EMail:                   email,
		Password:                pwd,
		EmailVerified:           false,
		TwoFactorAuthentication: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	username := StringWithCharset(30)
	updateUsernameAuthObj := Auth{
		EMail:    email,
		Username: username,
	}

	err = updateUsernameAuthObj.SetUsername()
	if err != nil {
		t.Errorf("Error while updating username: %s", err)
	}

	uName, err := auth.GetUsername()
	if err != nil {
		t.Errorf("Error while getting username: %s", err)
	}

	Equal(t, username, uName)
	Equal(t, nil, err)
}

func TestUpdateUsernameWithOver32Charcters(t *testing.T) {
	Setup()

	pwd := StringWithCharset(20)
	email := StringWithCharset(10) + "@gmail.com"
	auth := Auth{
		EMail:                   email,
		Password:                pwd,
		EmailVerified:           false,
		TwoFactorAuthentication: false,
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user: %s", err)
	}

	username := StringWithCharset(34)
	updateUsernameAuthObj := Auth{
		EMail:    email,
		Username: username,
	}

	err = updateUsernameAuthObj.SetUsername()
	if err == nil {
		t.Errorf("No Error thrown even though the username is over 32 characters")
	}

	Equal(t, "username can only be a maximum of 32 characters long", err.Error())
	NotEqual(t, nil, err)
}

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
