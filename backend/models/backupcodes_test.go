package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urento/shoppinglist/pkg/util"
)

func TestGenerateCodes(t *testing.T) {
	Setup()

	email := util.StringWithCharset(10) + "@gmail.com"

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

func TestGetCodes(t *testing.T) {
	Setup()

	email := util.StringWithCharset(10) + "@gmail.com"

	c, err := GenerateCodes(email, false)
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

	s1 := util.StringArrayToArray(codes, 0)
	s2 := util.StringArrayToArray(codes, 1)
	s3 := util.StringArrayToArray(codes, 2)
	s4 := util.StringArrayToArray(codes, 3)
	s5 := util.StringArrayToArray(codes, 4)

	Equal(t, c[0], s1)
	Equal(t, c[1], s2)
	Equal(t, c[2], s3)
	Equal(t, c[3], s4)
	Equal(t, c[4], s5)
	Equal(t, true, has)
	Equal(t, nil, err)
}

func TestVerifyCode(t *testing.T) {
	Setup()

	email := util.StringWithCharset(10) + "@gmail.com"

	c, err := GenerateCodes(email, false)
	if err != nil {
		t.Errorf("Error while generating codes: %s", err)
	}

	has, err := HasCodes(email)
	if err != nil {
		t.Errorf("Error while checking if the user already has backup codes: %s", err)
	}

	s := util.StringArrayToArray(c, 0)
	ok, err := VerifyCode(email, s)
	if err != nil {
		t.Errorf("Error while verifying the backup code: %s", err)
	}

	Equal(t, true, has)
	Equal(t, true, ok)
	Equal(t, nil, err)
}

func TestRemoveCodes(t *testing.T) {
	Setup()

	email := util.StringWithCharset(30) + "@gmail.com"

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
