package util

import (
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/cache"
)

func TestGenerateTokenAndParse(t *testing.T) {
	cache.Setup()

	t.Run("Generate Token and Parse", func(t *testing.T) {
		email := RandomString(10) + "@gmail.com"

		token, err := GenerateToken(email, false)
		tasdg := time.Now()
		if err != nil {
			t.Errorf("Error while generating token: %s", err)
		}

		parsed, err := ParseToken(token)
		if err != nil {
			t.Errorf("Error while parsing token: %s", err)
		}

		Equal(t, nil, err)
		Equal(t, parsed.ExpiresAt, tasdg.Add(24*time.Hour).Unix())
		Equal(t, email, parsed.Email)
		NotEqual(t, nil, token)
	})

	t.Run("Generate Refresh Token and Parse", func(t *testing.T) {
		email := RandomString(10) + "@gmail.com"

		token, err := GenerateToken(email, true)
		tadsg := time.Now()
		if err != nil {
			t.Errorf("Error while generating token: %s", err)
		}

		parsed, err := ParseToken(token)
		if err != nil {
			t.Errorf("Error while parsing token: %s", err)
		}

		Equal(t, nil, err)
		Equal(t, parsed.ExpiresAt, tadsg.Add(168*time.Hour).Unix())
		Equal(t, email, parsed.Email)
		NotEqual(t, nil, token)
	})
}
