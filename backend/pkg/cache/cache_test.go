package cache

import (
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

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(155)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	Setup()

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

	//reconnect to redis because redislab only allows 30 simultaneous connections and i close the redis connection after every request
	Setup()

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

	Setup()

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

	Setup()

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

	Setup()

	ok, err := DeleteTokenByEmail(email, token)
	if err != nil || !ok {
		t.Errorf("Error while deleting Token by Email: %s", err)
	}

	Setup()

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Error("Token still cached")
	}

	Setup()

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

	Setup()

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Errorf("No Error thrown 3")
	}

	Setup()

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

	Setup()

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

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
