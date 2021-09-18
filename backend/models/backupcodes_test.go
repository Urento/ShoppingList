package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestGenerateCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	t.Run("normal generate codes", func(t *testing.T) {
		_, err := GenerateCodes(email, false)
		if err != nil {
			t.Errorf("Error while generating codes: %s", err)
		}

		has, err := HasCodes(email)
		if err != nil {
			t.Errorf("Error while checking if the user already has backup codes: %s", err)
		}

		Equal(t, true, has)
		Equal(t, nil, err)
	})

	t.Run("regenerate codes", func(t *testing.T) {
		has2, err := HasCodes(email)
		if err != nil {
			t.Errorf("Error while checking if the user already has backup codes: %s", err)
		}

		_, err = GenerateCodes(email, true)
		if err != nil {
			t.Errorf("Error while generating codes: %s", err)
		}

		has, err := HasCodes(email)
		if err != nil {
			t.Errorf("Error while checking if the user already has backup codes: %s", err)
		}

		Equal(t, true, has2)
		Equal(t, true, has)
		Equal(t, nil, err)
	})
}

//TODO: ADD TEST
/*func TestGetCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	c, err := GenerateCodes(email, false)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}
	t.Log(c)

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
}*/

func TestVerifyCode(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	c, err := GenerateCodes(email, false)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}
	t.Log(c[0])

	ok, codes, err := VerifyCode(email, c[0])
	if err != nil {
		t.Errorf("Error while verifying the backup code: %s", err)
	}
	t.Log(codes)

	Equal(t, true, has)
	Equal(t, true, ok)
	Equal(t, nil, err)
}

func TestRemoveCodes(t *testing.T) {
	Setup()

	email := StringWithCharset(30) + "@gmail.com"

	_, err := GenerateCodes(email, false)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	hasBefore, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	err = RemoveCodes(email)
	if err != nil {
		t.Errorf("Error while deleting codes: %s", err)
	}

	hasAfter, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	Equal(t, true, hasBefore)
	Equal(t, true, hasAfter)
	Equal(t, nil, err)
}
