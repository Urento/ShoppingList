package models

import (
	"errors"
	"net/mail"

	"github.com/alexedwards/argon2id"
	"github.com/urento/shoppinglist/pkg/cache"
)

/**
* TODO: Cache Users and get user data from cache
 */

type Auth struct {
	ID            int    `gorm:"primary_key" json:"id"`
	EMail         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Rank          string `json:"rank"` //admin or default
}

func GetPasswordHash(email string) (string, error) {
	var password string
	err := db.Model(&Auth{}).Select("password").Where("e_mail = ?", email).First(&password).Error
	if err != nil {
		return "", err
	}
	return password, nil
}

func CheckAuth(email, password string) (bool, error) {
	exists, err := Exists(email)
	if err != nil || !exists {
		return false, nil
	}

	var auth Auth
	pwdHash, err1 := GetPasswordHash(email)
	if err1 != nil {
		return false, nil
	}

	match, err := argon2id.ComparePasswordAndHash(password, pwdHash)
	if err != nil {
		return false, nil
	}

	authEmail := Auth{
		EMail:    email,
		Password: pwdHash,
	}

	err = db.Select("id").Model(&Auth{}).Where(authEmail).First(&auth).Error
	if err != nil {
		return false, err
	}
	if auth.ID > 0 && match {
		return true, nil
	}

	return false, nil
}

func GetUser(email string) (*Auth, error) {
	var user Auth
	err := db.Model(&Auth{}).Where("e_mail = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateAccount(email, username, password string) error {
	validEmail := validateEmail(email)
	if !validEmail {
		return errors.New("email is not valid")
	}

	exists, err := Exists(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("account already exists")
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	authObj := Auth{
		EMail:         email,
		Username:      username,
		Password:      passwordHash,
		EmailVerified: false,
		Rank:          "default",
	}

	err = db.Create(&authObj).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteAccount(email, password string) error {
	pwdHash, err := GetPasswordHash(email)
	if err != nil {
		return err
	}

	match, err := argon2id.ComparePasswordAndHash(password, pwdHash)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("wrong password")
	}

	err = db.Where("e_mail = ?", email).Delete(&Auth{}).Error
	if err != nil {
		return err
	}
	return nil
}

func Logout(email, token string) (bool, error) {
	ok, err := cache.DeleteTokenByEmail(email, token)
	if err != nil || !ok {
		return false, err
	}
	return true, nil
}

func GetRank(email string) (string, error) {
	var rank string
	if err := db.Model(&Auth{}).Where("e_mail = ?", email).Select("rank").First(&rank).Error; err != nil {
		return "", err
	}
	return rank, nil
}

func SetRank(email, rank string) error {
	rankCheck := rankExists(rank)
	if !rankCheck {
		return errors.New("rank does not exist")
	}

	if err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("rank", rank).Error; err != nil {
		return err
	}
	return nil
}

func rankExists(rank string) bool {
	if rank == "default" || rank == "admin" {
		return true
	}
	return false
}

func Exists(email string) (exists bool, err error) {
	var Found bool
	err = db.Raw("SELECT EXISTS(SELECT id FROM auths WHERE e_mail = ?) AS found", email).Scan(&Found).Error
	return Found, err
}

func IsEmailVerified(email string) (bool, error) {
	var verified bool
	err := db.Model(&Auth{}).Select("email_verified").Where("e_mail = ?", email).First(&verified).Error
	if err != nil {
		return false, err
	}
	if !verified {
		return false, nil
	}
	return true, nil
}

func VerifyEmail(email string) error {
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("email_verified", true).Error
	if err != nil {
		return err
	}
	return nil
}

func SendVerifyEmail(email string) error {
	//TODO: SEND EMAIL
	return errors.New("not implemented yet")
}

func Count(email string) (int64, error) {
	count := int64(0)
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Count(&count).Error
	if err != nil {
		return 1000, err
	}
	return count, nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
