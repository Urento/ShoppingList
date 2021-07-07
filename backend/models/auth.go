package models

import (
	"github.com/alexedwards/argon2id"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	EMail    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetPasswordHash(email, username string) (string, error) {
	var password string
	if email != "" {
		err := db.Select("password").Where("email = ?", email).First(&password).Error
		if err != nil {
			return "", err
		}
	} else {
		err := db.Select("password").Where("username = ?", username).First(&password).Error
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func CheckAuth(email, username, password string) (bool, error) {
	var auth Auth
	pwdHash, err1 := GetPasswordHash(email, username)
	if err1 != nil {
		return false, nil
	}

	match, err := argon2id.ComparePasswordAndHash(password, pwdHash)
	if err != nil {
		return false, nil
	}

	// auth with username
	if username != "" {
		err := db.Select("id").Where(&Auth{Username: username, Password: pwdHash}).First(&auth).Error
		if err != nil {
			return false, err
		}
		if auth.ID > 0 && match {
			return true, nil
		}

		return false, nil
	}

	// auth with email
	err = db.Select("id").Where(&Auth{EMail: username, Password: pwdHash}).First(&auth).Error
	if err != nil {
		return false, err
	}
	if auth.ID > 0 && match {
		return true, nil
	}

	return false, nil
}

func CreateAccount(email, username, password string) error {
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = db.Model(&Auth{}).Create(Auth{EMail: email, Username: username, Password: passwordHash}).Error
	if err != nil {
		return err
	}
	return nil
}
