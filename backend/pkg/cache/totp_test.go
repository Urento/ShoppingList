package cache

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCacheAndGetTOTPSecret(t *testing.T) {
	Setup(false)

	t.Run("TestCacheAndGetTOTPSecret", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		secret := StringWithCharset(100)

		err := CacheTOTPSecret(email, secret)
		if err != nil {
			t.Errorf("Error while caching TOTP Secret: %s", err)
		}

		totpSecret, err := GetTOTPSecret(email)
		if err != nil {
			t.Errorf("Error while getting TOTP Secret: %s", err)
		}

		Equal(t, secret, totpSecret)
		Equal(t, nil, err)
	})

	t.Run("TestGetTOTPSecretThatDoesntExist", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		_, err := GetTOTPSecret(email)

		Equal(t, "totp secret is not cached", err.Error())
	})
}

func TestDeleteTOTPSecret(t *testing.T) {
	Setup(false)

	email := StringWithCharset(100) + "@gmail.com"
	secret := StringWithCharset(100)

	err := CacheTOTPSecret(email, secret)
	if err != nil {
		t.Errorf("Error while caching TOTP Secret: %s", err)
	}

	totpSecret, err := GetTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while getting TOTP Secret: %s", err)
	}

	err = DeleteTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while deleting TOTP Secret: %s", err)
	}

	_, delErr := GetTOTPSecret(email)

	Equal(t, secret, totpSecret)
	Equal(t, nil, err)
	Equal(t, "totp secret is not cached", delErr.Error())
}

func TestIsTOTPCached(t *testing.T) {
	Setup(false)

	t.Run("TestIsTOTPCached", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		secret := StringWithCharset(100)

		err := CacheTOTPSecret(email, secret)
		if err != nil {
			t.Errorf("Error while caching TOTP Secret: %s", err)
		}

		totpSecret, err := GetTOTPSecret(email)
		if err != nil {
			t.Errorf("Error while getting TOTP Secret: %s", err)
		}

		ok, err := IsTOTPSecretCached(email)
		if err != nil {
			t.Errorf("Error while checking if TOTP Secret is cached: %s", err)
		}

		Equal(t, secret, totpSecret)
		Equal(t, true, ok)
	})

	t.Run("TestIsTOTPCachedWhenItsNotCached", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		_, err := GetTOTPSecret(email)

		ok, _ := IsTOTPSecretCached(email)

		Equal(t, false, ok)
		Equal(t, "totp secret is not cached", err.Error())
	})
}
