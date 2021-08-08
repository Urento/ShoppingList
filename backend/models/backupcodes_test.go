package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestGenerateCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	_, err := GenerateCodes(email)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	Equal(t, true, has)
	Equal(t, nil, err)
}

func TestGenerateCodesWhenTheUserAlreadyHasCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	codes, err := GenerateCodes(email)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	codes2, err2 := GenerateCodes(email)
	if err2 != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has2, err2 := HasCodes(email)
	if err2 != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	Equal(t, true, has)
	Equal(t, true, has2)
	Equal(t, nil, err)
	NotEqual(t, codes2, codes)
}

func TestGetCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	c, err := GenerateCodes(email)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	codes, err := GetCodes(email)
	if err != nil {
		t.Errorf("Error while getting codes: %s", err)
	}
	t.Log(codes)

	Equal(t, c, codes)
	Equal(t, true, has)
	Equal(t, nil, err)
}
