package cache

import (
	"testing"

	"github.com/alexedwards/argon2id"
	. "github.com/stretchr/testify/assert"
)

func TestCacheUser(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"
	username := StringWithCharset(100)
	password := StringWithCharset(100)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

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
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	Equal(t, nil, err)
}

func TestGetUser(t *testing.T) {
	Setup()

	t.Run("TestGetUser", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		username := StringWithCharset(100)
		password := StringWithCharset(100)
		emailVerified := RandomBoolean()
		rank := RandomRank()
		twoFactorAuthentication := RandomBoolean()
		ip := RandomIPAddress()

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
			IPAddress:               ip,
		}

		err = u.CacheUser()
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
	})

	t.Run("TestGetUserThatDoesntExist", func(t *testing.T) {
		_, err := GetUser("dkjfgbksdjhfgbkjdhfsgb@gmail.com")

		NotEqual(t, nil, err)
	})
}

func TestUpdateUser(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"
	username := StringWithCharset(100)
	password := StringWithCharset(100)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

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
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	newUsername := StringWithCharset(100)
	newEmailVerified := RandomBoolean()
	newRank := RandomRank()
	newTwoFactorAuthentication := RandomBoolean()
	newIp := RandomIPAddress()
	newUser := User{
		EMail:                   email,
		Username:                newUsername,
		Password:                pwdHash,
		EmailVerified:           newEmailVerified,
		Rank:                    newRank,
		TwoFactorAuthentication: newTwoFactorAuthentication,
		IPAddress:               newIp,
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

	email := StringWithCharset(100) + "@gmail.com"
	username := StringWithCharset(100)
	password := StringWithCharset(100)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

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
		IPAddress:               ip,
	}

	err = u.CacheUser()
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
