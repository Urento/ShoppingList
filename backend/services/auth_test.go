package services

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	Setup()

	auth := Auth{
		EMail:    StringWithCharset(10) + "@gmail.com",
		Username: StringWithCharset(10),
		Password: StringWithCharset(20),
	}

	err := auth.Create()
	if err != nil {
		t.Errorf("Error while creating user %s", err.Error())
	}

	check, err := auth.Check()
	if err != nil {
		t.Errorf("Error while checking user %s", err.Error())
	}

	if !check {
		t.Errorf("check was false")
	}

	Equal(t, check, true)
	Equal(t, err, nil)
}
