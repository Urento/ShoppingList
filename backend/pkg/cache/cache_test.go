package cache

import (
	"math/rand"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	. "github.com/stretchr/testify/assert"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestCacheJWTToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)

	Equal(t, nil, err)
	Equal(t, true, exists)
}

func TestGetTokenByEmail(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	val, err := GetJWTByEmail(email)
	if err != nil {
		t.Errorf("Error while getting Token by Email %s", err)
	}

	Equal(t, nil, err)
	Equal(t, token, val)
}

func TestDoesTokenExpireAfter1Day(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	ttl, err := GetTTLByEmail(email)
	if err != nil {
		t.Errorf("Error getting the ttl from the key by email %s", err)
	}

	if ttl < 86200 {
		t.Errorf("ttl is too low")
	}

	Equal(t, nil, err)
}

func TestGetEmailByJWT(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	val, err := GetEmailByJWT(token)
	if err != nil {
		t.Errorf("Error while getting Email by Token: %s", err)
	}

	Equal(t, email, val)
}

func TestDeleteToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	ok, err := DeleteTokenByEmail(email, token)
	if err != nil || !ok {
		t.Errorf("Error while deleting Token by Email: %s", err)
	}

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Error("Token still cached")
	}

	_, err = GetEmailByJWT(token)
	if err == nil {
		t.Error("Token is still cached")
	}

	Equal(t, true, ok)
	Equal(t, "jwt token not cached", err.Error())
}

func TestDeleteTokenWithEmailThatDoesntExist(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	ok, err := DeleteTokenByEmail(email, token)
	if err == nil || ok {
		t.Errorf("No Error thrown 4")
	}

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Errorf("No Error thrown 3")
	}

	_, err = GetEmailByJWT(token)
	if err == nil {
		t.Errorf("No Error thrown 2 ")
	}

	Equal(t, false, ok)
	Equal(t, "jwt token not cached", err.Error())
}

func TestIsTokenValidWithValidToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	valid, err := IsTokenValid(token)
	if err != nil {
		t.Errorf("Error while checking if token is valid %s", err)
	}

	Equal(t, true, valid)
	Equal(t, nil, err)
}

func TestIsTokenValidWithInvalidToken(t *testing.T) {
	Setup()

	token := StringWithCharset(155)

	valid, _ := IsTokenValid(token)

	Equal(t, false, valid)
}

func TestCacheUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
	}

	err = CacheUser(u)
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	Equal(t, nil, err)
}

func TestGetUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
	}

	err = CacheUser(u)
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	Equal(t, nil, err)
	Equal(t, email, user.EMail)
	Equal(t, pwdHash, user.Password)
	Equal(t, emailVerified, user.EmailVerified)
	Equal(t, username, user.Username)
	Equal(t, rank, user.Rank)
	Equal(t, twoFactorAuthentication, user.TwoFactorAuthentication)
}

func TestGetUserThatDoesntExist(t *testing.T) {
	Setup()

	_, err := GetUser("dkjfgbksdjhfgbkjdhfsgb@gmail.com")

	NotEqual(t, nil, err)
}

func TestUpdateUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
	}

	err = CacheUser(u)
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	newUsername := StringWithCharset(10)
	newEmailVerified := RandomBoolean()
	newRank := RandomRank()
	newTwoFactorAuthentication := RandomBoolean()
	newUser := User{
		EMail:                   email,
		Username:                newUsername,
		Password:                pwdHash,
		EmailVerified:           newEmailVerified,
		Rank:                    newRank,
		TwoFactorAuthentication: newTwoFactorAuthentication,
	}

	err = UpdateUser(newUser)
	if err != nil {
		t.Errorf("Error while updating user: %s", err)
	}

	updatedUser, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting updated user: %s", err)
	}

	Equal(t, nil, err)
	Equal(t, email, user.EMail)
	Equal(t, pwdHash, user.Password)
	Equal(t, emailVerified, user.EmailVerified)
	Equal(t, username, user.Username)
	Equal(t, rank, user.Rank)
	Equal(t, twoFactorAuthentication, user.TwoFactorAuthentication)
	Equal(t, email, updatedUser.EMail)
	Equal(t, newUsername, updatedUser.Username)
	Equal(t, pwdHash, updatedUser.Password)
	Equal(t, newEmailVerified, updatedUser.EmailVerified)
	Equal(t, newRank, updatedUser.Rank)
	Equal(t, newTwoFactorAuthentication, updatedUser.TwoFactorAuthentication)
}

func TestDeleteUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
	}

	err = CacheUser(u)
	if err != nil {
		t.Errorf("Error while caching user: %s", err)
	}

	err = DeleteUser(email)
	if err != nil {
		t.Errorf("Error while deleting user: %s", err)
	}

	_, shouldErr := GetUser(email)
	if shouldErr == nil {
		t.Errorf("GetUser didn't throw an error after deleting")
	}

	Equal(t, nil, err)
	NotEqual(t, nil, shouldErr)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomBoolean() bool {
	nmb := rand.Intn(2)
	return nmb <= 1
}

func RandomRank() string {
	nmb := rand.Intn(2)
	if nmb <= 1 {
		return "admin"
	}
	return "default"
}
