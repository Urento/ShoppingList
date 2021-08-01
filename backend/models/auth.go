package models

import (
	"errors"
	"net/mail"

	"github.com/alexedwards/argon2id"
	"github.com/urento/shoppinglist/pkg/cache"
	"gorm.io/gorm"
)

/**
* TODO: Cache Users and get user data from cache
 */

type Auth struct {
	Model

	ID                      int    `gorm:"primary_key" json:"id"`
	EMail                   string `json:"e_mail"`
	EmailVerified           bool   `json:"email_verified"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	Rank                    string `json:"rank"` //admin or default
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
	IPAddress               string `json:"ip_address"`
}

func GetPasswordHash(email string) (string, error) {
	var password string
	err := db.Model(&Auth{}).Select("password").Where("e_mail = ?", email).First(&password).Error
	return password, err
}

func CheckAuth(email, password, ip string) (bool, error) {
	exists, err := Exists(email)
	if err != nil || !exists {
		return false, nil
	}

	pwdHash, err1 := GetPasswordHash(email)
	if err1 != nil {
		return false, nil
	}

	match, err := argon2id.ComparePasswordAndHash(password, pwdHash)
	if err != nil {
		return false, nil
	}

	var auth Auth
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Select("id").Model(&Auth{}).Where("e_mail = ?", email).First(&auth).Error
		err = tx.Model(&Auth{}).Where("e_mail = ?", email).Update("ip_address", ip).Error
		return err
	})

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
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Select("id, e_mail, email_verified, username, rank, two_factor_authentication, created_on, modified_on, deleted_at").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateAccount(email, username, password, ip string) error {
	validEmail := validateEmail(email)
	if !validEmail {
		return errors.New("email is not valid")
	}

	exists, err := Exists(email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email is already being used")
	}

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	authObj := Auth{
		EMail:                   email,
		Password:                passwordHash,
		Username:                username,
		EmailVerified:           false,
		Rank:                    "default",
		TwoFactorAuthentication: false,
		IPAddress:               ip,
	}

	err = db.Create(&authObj).Error
	if err != nil {
		return err
	}

	return nil
}

func SetUsername(email, username string) error {
	if len(username) > 32 {
		return errors.New("username can only be a maximum of 32 characters long")
	}

	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("username", username).Error
	return err
}

func GetUsername(email string) (string, error) {
	var username string
	err := db.Model(&Auth{}).Select("username").Where("e_mail = ?", email).Find(&username).Error
	return username, err
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
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Select("rank").First(&rank).Error
	return rank, err
}

func SetRank(email, rank string) error {
	rankCheck := rankExists(rank)
	if !rankCheck {
		return errors.New("rank does not exist")
	}

	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("rank", rank).Error

	return err
}

func rankExists(rank string) bool {
	return rank == "default" || rank == "admin"
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
	return verified, nil
}

func VerifyEmail(email string) error {
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("email_verified", true).Error
	return err
}

func SendVerifyEmail(email string) error {
	//TODO: SEND EMAIL
	return errors.New("not implemented yet")
}

func SetTwoFactorAuthentication(email string, status bool) error {
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("two_factor_authentication", status).Error
	return err
}

func IsTwoFactorEnabled(email string) (bool, error) {
	var status bool
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Select("two_factor_authentication").First(&status).Error
	if err != nil {
		return false, err
	}

	return status, nil
}

func UpdateIP(email, ip string) error {
	err := db.Model(&Auth{}).Where("e_mail = ?", email).Update("ip_address", ip).Error
	return err
}

func GetIP(email string) (string, error) {
	var ip string
	err := db.Model(&Auth{}).Select("ip_address").Where("e_mail = ?", email).Find(&ip).Error
	return ip, err
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
