package util

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestGenerateTokenAndParse(t *testing.T) {
	cache.Setup()

	email := StringWithCharset(10) + "@gmail.com"

	token, err := GenerateToken(email)
	if err != nil {
		t.Errorf("Error while generating token: %s", err)
	}

	parsed, err := ParseToken(token)
	if err != nil {
		t.Errorf("Error while parsing token: %s", err)
	}

	Equal(t, nil, err)
	Equal(t, email, parsed.Email)
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
