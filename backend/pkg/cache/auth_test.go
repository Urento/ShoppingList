package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

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

	Nil(t, err)
	True(t, exists)
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

	Nil(t, err)
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

	Nil(t, err)
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

		True(t, ok)
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

		False(t, ok)
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

		True(t, valid)
		Nil(t, err)
	})

	t.Run("TestIsTokenValidWithInvalidToken", func(t *testing.T) {
		token := StringWithCharset(245)

		valid, _ := IsTokenValid(token)

		False(t, valid)
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

		True(t, ok)
		Nil(t, err)
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

		False(t, ok)
		Nil(t, err)
	})

	t.Run("TestVerifySecretIdWithWrongIdWithoutAccount", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		ok, err := VerifySecretId(email, "secretId")
		if err == nil {
			t.Errorf("No error thrown even though the secretId is wrong and doesn't exist")
		}

		False(t, ok)
		Equal(t, "secretid is not valid", err.Error())
	})
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

		True(t, ok)
		Nil(t, err)
		Equal(t, secretId, key)
	})

	t.Run("TestHasSecretIdWhenTheUserDoesntHaveOne", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"

		_, has, _ := HasSecretId(email)

		if has {
			t.Errorf("User doesnt have a SecretId but it says it has")
		}

		False(t, has)
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

	True(t, ok)
	Nil(t, err)
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

		True(t, exists)
		True(t, ok)
		False(t, ok2)
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

	True(t, exists)
	True(t, ok)
}

func TestGetFailedLoginAttemts(t *testing.T) {
	Setup()

	t.Run("Get Failed Login Attempts", func(t *testing.T) {
		email := StringWithCharset(500)
		ctx := context.Background()

		err := UpdateFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while updating failed login attempts: %s", err)
		}

		failedAttemts, err := GetFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while getting failed login attempts: %s", err)
		}

		Equal(t, 1, failedAttemts)
		Nil(t, err)
	})

	t.Run("Get Failed Login Attempts with multiple attempts", func(t *testing.T) {
		email := StringWithCharset(500)
		ctx := context.Background()

		for i := 1; i < 10; i++ {
			err := UpdateFailedLoginAttempts(ctx, email)
			if err != nil {
				t.Errorf("Error while updating failed login attempts: %s", err)
			}
		}

		failedAttemts, err := GetFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while getting failed login attempts: %s", err)
		}

		Equal(t, 9, failedAttemts)
		Nil(t, err)
	})
}

func TestHasFailedLoginAttempts(t *testing.T) {
	Setup()

	t.Run("Has Failed Login Attempts", func(t *testing.T) {
		email := StringWithCharset(100) + "@gmail.com"
		ctx := context.Background()

		hasBefore, err := HasFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while checking if the user has failed login attempts: %s", err)
		}

		err = UpdateFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while updating failed login attempts: %s", err)
		}

		has, err := HasFailedLoginAttempts(ctx, email)
		if err != nil {
			t.Errorf("Error while checking if the user has failed login attempts: %s", err)
		}

		True(t, has)
		False(t, hasBefore)
	})

	t.Run("Has Failed Login Attempts when the user doesn't exist", func(t *testing.T) {
		has, _ := HasFailedLoginAttempts(context.Background(), "dfkgnj")

		False(t, has)
	})
}

func TestClearFailedLoginAttempts(t *testing.T) {
	Setup()

	email := StringWithCharset(100) + "@gmail.com"
	ctx := context.Background()

	err := UpdateFailedLoginAttempts(ctx, email)
	if err != nil {
		t.Errorf("Error while updating failed login attempts: %s", err)
	}

	has, err := HasFailedLoginAttempts(ctx, email)
	if err != nil {
		t.Errorf("Error while checking if the user has failed login attempts: %s", err)
	}

	err = ClearFailedLoginAttempts(ctx, email)
	if err != nil {
		t.Errorf("Error while clearing failed login attempts: %s", err)
	}

	attempts, err := GetFailedLoginAttempts(ctx, email)
	if err != nil {
		t.Errorf("Error while getting failed login attempts: %s", err)
	}

	True(t, has)
	Equal(t, 0, attempts)
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
