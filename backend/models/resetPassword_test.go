package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCreateResetPassword(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	err := CreateResetPassword(email)
	if err != nil {
		t.Errorf("Error while creating reset password: %s", err)
	}

	exists, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request exists: %s", err)
	}

	Equal(t, true, exists)
	Equal(t, nil, err)
}

func TestExistResetPasswordWhereRequestDoesNotExist(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	exists, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password object exists: %s", err)
	}

	Equal(t, false, exists)
	Equal(t, nil, err)
}

func TestDeleteResetPassword(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	err := CreateResetPassword(email)
	if err != nil {
		t.Errorf("Error while creating reset password: %s", err)
	}

	existsBefore, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request exists 1: %s", err)
	}

	err = DeleteResetPassword(email)
	if err != nil {
		t.Errorf("Error while deleting reset password object: %s", err)
	}

	existsAfter, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request exists 1: %s", err)
	}

	Equal(t, true, existsBefore)
	Equal(t, false, existsAfter)
	Equal(t, nil, err)
}

func TestDeleteResetPasswordWhereRequestDoesNotExist(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	existsBefore, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request exists 1: %s", err)
	}

	err = DeleteResetPassword(email)
	if err != nil {
		t.Errorf("Error while deleting reset password object: %s", err)
	}

	existsAfter, err := ExistsResetPassword(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request exists 1: %s", err)
	}

	Equal(t, false, existsBefore)
	Equal(t, false, existsAfter)
	Equal(t, nil, err)
}

func TestVerifyVerificationId(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	err := CreateResetPassword(email)
	if err != nil {
		t.Errorf("Error while creating reset password: %s", err)
	}

	verificationId, err := GetVerificationId(email)
	if err != nil {
		t.Errorf("Error while getting verification id: %s", err)
	}

	ok, err := VerifyVerificationId(email, verificationId)
	if err != nil {
		t.Errorf("Error while verifying verification id: %s", err)
	}

	Equal(t, true, ok)
	Equal(t, nil, err)
}

func TestVerifyVerificationIdWithWrongId(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	err := CreateResetPassword(email)
	if err != nil {
		t.Errorf("Error while creating reset password: %s", err)
	}

	ok, err := VerifyVerificationId(email, "verificationId")

	Equal(t, false, ok)
	Equal(t, nil, err)
}
