package models

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCreateResetPassword(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	t.Run("TestCreateResetPassword", func(t *testing.T) {
		err := CreateResetPassword(email)
		if err != nil {
			t.Errorf("Error while creating reset password: %s", err)
		}

		exists, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists: %s", err)
		}

		Equal(t, true, exists)
		Equal(t, nil, err)
	})

	t.Run("TestCreateResetPasswordWhenRequestAlreadyExists", func(t *testing.T) {
		existsBefore, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists: %s", err)
		}

		verificationIDBefore, err := GetVerificationID(email)
		if err != nil {
			t.Errorf("Error while getting verificationID: %s", err)
		}

		err = CreateResetPassword(email)
		if err != nil {
			t.Errorf("Error while creating reset password: %s", err)
		}

		existsAfter, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists: %s", err)
		}

		verificationIDAfter, err := GetVerificationID(email)
		if err != nil {
			t.Errorf("Error while getting verificaitonID: %s", err)
		}

		Equal(t, true, existsBefore)
		Equal(t, true, existsAfter)
		Equal(t, nil, err)
		Equal(t, verificationIDBefore, verificationIDAfter)
	})
}

func TestExistResetPassword(t *testing.T) {
	Setup()

	t.Run("Exist Reset Password Where Request doesn't exist", func(t *testing.T) {
		email := StringWithCharset(10) + "@gmail.com"

		exists, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password object exists: %s", err)
		}

		Equal(t, false, exists)
		Equal(t, nil, err)
	})
}

func TestDeleteResetPassword(t *testing.T) {
	Setup()

	t.Run("Delete Reset Password", func(t *testing.T) {
		email := StringWithCharset(10) + "@gmail.com"

		err := CreateResetPassword(email)
		if err != nil {
			t.Errorf("Error while creating reset password: %s", err)
		}

		existsBefore, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists 1: %s", err)
		}

		err = DeleteResetPassword(email)
		if err != nil {
			t.Errorf("Error while deleting reset password object: %s", err)
		}

		existsAfter, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists 1: %s", err)
		}

		Equal(t, true, existsBefore)
		Equal(t, false, existsAfter)
		Equal(t, nil, err)
	})

	t.Run("Delete Reset Password Where Request Doesn't exist", func(t *testing.T) {
		email := StringWithCharset(10) + "@gmail.com"

		existsBefore, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists 1: %s", err)
		}

		err = DeleteResetPassword(email)
		if err != nil {
			t.Errorf("Error while deleting reset password object: %s", err)
		}

		existsAfter, err := HasResetPassword(email)
		if err != nil {
			t.Errorf("Error while checking if reset password request exists 1: %s", err)
		}

		Equal(t, false, existsBefore)
		Equal(t, false, existsAfter)
		Equal(t, nil, err)
	})
}

func TestVerifyVerificationId(t *testing.T) {
	Setup()

	t.Run("Verify Verification Id", func(t *testing.T) {
		email := StringWithCharset(10) + "@gmail.com"

		err := CreateResetPassword(email)
		if err != nil {
			t.Errorf("Error while creating reset password: %s", err)
		}

		verificationId, err := GetVerificationID(email)
		if err != nil {
			t.Errorf("Error while getting verification id: %s", err)
		}

		ok, err := VerifyVerificationID(email, verificationId)
		if err != nil {
			t.Errorf("Error while verifying verification id: %s", err)
		}

		Equal(t, true, ok)
		Equal(t, nil, err)
	})

	t.Run("Test Verify Verification Id With Wrong Id", func(t *testing.T) {
		email := StringWithCharset(10) + "@gmail.com"

		err := CreateResetPassword(email)
		if err != nil {
			t.Errorf("Error while creating reset password: %s", err)
		}

		ok, err := VerifyVerificationID(email, "verificationId")

		Equal(t, false, ok)
		Equal(t, nil, err)
	})
}

func TestIsStillValid(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	err := CreateResetPassword(email)
	if err != nil {
		t.Errorf("Error while creating reset password: %s", err)
	}

	valid, err := IsStillValid(email)
	if err != nil {
		t.Errorf("Error while checking if reset password request is still valid: %s", err)
	}

	Equal(t, true, valid)
	Equal(t, nil, err)
}
