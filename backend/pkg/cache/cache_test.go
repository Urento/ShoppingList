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
	t.Log(email)

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

	t.Log(token)

	Equal(t, email, val)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
