package util

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestIsEmailValid(t *testing.T) {
	validEmail := RandomEmail()
	invalidEmail := StringWithCharset(100)
	t.Log(validEmail)

	isValid1 := IsEmailValid(validEmail)
	isValid2 := IsEmailValid(invalidEmail)

	True(t, isValid1)
	False(t, isValid2)
}
