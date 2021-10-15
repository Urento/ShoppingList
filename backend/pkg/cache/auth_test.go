package cache

import (
	"fmt"
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

	email := StringWithCharset(100) + "@gmail.com"
	token := StringWithCharset(245)

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

	email := StringWithCharset(100) + "@gmail.com"
	token := StringWithCharset(245)

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

	email := StringWithCharset(100) + "@gmail.com"
	token := StringWithCharset(245)

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

	email := StringWithCharset(100) + "@gmail.com"
	token := StringWithCharset(245)

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

	t.Run("TestDeleteToken", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		token := StringWithCharset(245)

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
	})

	t.Run("TestDeleteTokenWithEmailThatDoesntExist", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		token := StringWithCharset(245)

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
	})
}

func TestIsTokenValid(t *testing.T) {
	Setup()

	t.Run("TestIsTokenValid", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		token := StringWithCharset(245)

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
	})

	t.Run("TestIsTokenValidWithInvalidToken", func(t *testing.T) {
		token := StringWithCharset(245)

		valid, _ := IsTokenValid(token)

		Equal(t, false, valid)
	})
}

func TestGenerateSecretIdAndVerify(t *testing.T) {
	Setup()

	t.Run("TestGenerateSecretIdAndVerify", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		secretId, err := GenerateSecretId(email)
		if err != nil {
			t.Errorf("Error while generating secret id: %s", err)
		}

		ok, err := VerifySecretId(email, secretId)
		if err != nil {
			t.Errorf("Error while verifying secert id: %s", err)
		}

		Equal(t, true, ok)
		Equal(t, nil, err)
	})

	t.Run("TestVerifySecretIdWithWrongIdWithExistingAccount", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		_, err := GenerateSecretId(email)
		if err != nil {
			t.Errorf("Error while generating secret id: %s", err)
		}

		ok, err := VerifySecretId(email, "secretId")
		if err != nil {
			t.Errorf("Error while verifying secert id: %s", err)
		}

		Equal(t, false, ok)
		Equal(t, nil, err)
	})

	t.Run("TestVerifySecretIdWithWrongIdWithoutAccount", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		ok, err := VerifySecretId(email, "secretId")
		if err == nil {
			t.Errorf("No error thrown even though the secretId is wrong and doesn't exist")
		}

		Equal(t, false, ok)
		Equal(t, "secretid is not valid", err.Error())
	})
}

func TestGetTwoFactorAuthenticationStatus(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"
	username := StringWithCharset(100)
	password := StringWithCharset(300)
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

	status, err := GetTwoFactorAuthenticationStatus(email)
	if err != nil {
		t.Errorf("Error while getting two factor authentication status: %s", err)
	}

	Equal(t, twoFactorAuthentication, status)
	Equal(t, nil, err)
}

func TestHasSecretId(t *testing.T) {
	Setup()

	t.Run("TestHasSecretId", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		secretId, err := GenerateSecretId(email)
		if err != nil {
			t.Errorf("Error while generating secret id: %s", err)
		}

		ok, err := VerifySecretId(email, secretId)
		if err != nil {
			t.Errorf("Error while verifying secert id: %s", err)
		}

		key, has, err := HasSecretId(email)
		if err != nil {
			t.Errorf("Error while checking if the user still has a secretId: %s", err)
		}

		if !has {
			t.Errorf("User does not have a secretId even though he has one")
		}

		if key != secretId {
			t.Errorf("SecretId is not the same as the previously generated one")
		}

		Equal(t, true, ok)
		Equal(t, nil, err)
		Equal(t, secretId, key)
	})

	t.Run("TestHasSecretIdWhenTheUserDoesntHaveOne", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		_, has, _ := HasSecretId(email)

		if has {
			t.Errorf("User doesnt have a SecretId but it says it has")
		}

		Equal(t, false, has)
	})
}

func TestInvalidateSecretId(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"

	secretId, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, secretId)
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	key, has, err := HasSecretId(email)
	if err != nil {
		t.Errorf("Error while checking if the user still has a secretId: %s", err)
	}

	if !has {
		t.Errorf("User does not have a secretId even though he has one")
	}

	if key != secretId {
		t.Errorf("SecretId is not the same as the previously generated one")
	}

	err = InvalidateSecretId(email)
	if err != nil {
		t.Errorf("Error while invalidating secretId: %s", err)
	}

	_, has2, err := HasSecretId(email)
	if err != nil {
		t.Errorf("Error while checking if the user still has a secretId 2: %s", err)
	}

	if has2 {
		t.Errorf("SecretId did not get invalidated!")
	}

	Equal(t, true, ok)
	Equal(t, nil, err)
	Equal(t, secretId, key)
}

func TestInvalidateJWTTokens(t *testing.T) {
	Setup()

	t.Run("Invalidate specific JWT Token", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		token := StringWithCharset(245)

		err := CacheJWT(email, token)
		if err != nil {
			t.Errorf("Error while caching JWT token %s", err)
		}

		exists, err := EmailExists(email)
		if err != nil {
			t.Errorf("Error while checking if the email exists: %s", err)
		}

		ok, err := DoesTokenBelongToEmail(email, token)
		if err != nil {
			t.Errorf("Error while checking if the token belongs to the email: %s", err)
		}

		err = InvalidateSpecificJWTToken(email, token)
		if err != nil {
			t.Errorf("Error while invalidating specific jwt token: %s", err)
		}

		ok2, _ := DoesTokenBelongToEmail(email, token)

		Equal(t, true, exists)
		Equal(t, true, ok)
		Equal(t, false, ok2)
	})
}

func TestDoesTokenBelongToEmail(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)
	if err != nil {
		t.Errorf("Error while checking if the email exists: %s", err)
	}

	ok, err := DoesTokenBelongToEmail(email, token)
	if err != nil {
		t.Errorf("Error while checking if the token belongs to the email: %s", err)
	}

	Equal(t, true, exists)
	Equal(t, true, ok)
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

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
