package util

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
)

func TestGenerateTokenAndParse(t *testing.T) {
	cache.Setup()

	email := RandomString(10) + "@gmail.com"

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
	NotEqual(t, nil, token)
}
