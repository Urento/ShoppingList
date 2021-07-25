package util

import (
	"math/rand"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestGenerateTokenAndParse(t *testing.T) {
	email := StringWithCharset(10) + "@gmail.com"
	password := StringWithCharset(30)

	hashedPwd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while hashing password %s", err)
	}

	cache.Setup()

	token, err := GenerateToken(email, hashedPwd)
	if err != nil {
		t.Errorf("Error while generating token %s", err)
	}

	parsed, err := ParseToken(token)
	if err != nil {
		t.Errorf("Error while parsing token %s", err)
	}

	pwdHashOk, err := argon2id.ComparePasswordAndHash(password, parsed.Password)

	Equal(t, nil, err)
	Equal(t, email, parsed.Email)
	Equal(t, true, pwdHashOk)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
