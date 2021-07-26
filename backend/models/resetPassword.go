package models

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ResetPassword struct {
	Model

	VerificationID  string `json:"verification_id"`
	Email           string `json:"email"`
	AlreadyVerified bool   `json:"already_verified"`
}

//Maybe change to redis?
//TODO: Let them expire after 1 hour (just check createdat and add 1 hour)

func ExistsResetPassword(email string) (bool, error) {
	var Exists bool
	err := db.Raw("SELECT EXISTS(SELECT created_on FROM reset_passwords WHERE email = ? AND already_verified = ?) AS found", email, false).Scan(&Exists).Error
	return Exists, err
}

func DeleteResetPassword(email string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Model(&ResetPassword{}).Where("email = ?", email).Update("already_verified", true).Error
		if err != nil {
			return err
		}

		err = db.Where("email = ?", email).Delete(&ResetPassword{}).Error
		return err
	})
	return err
}

func GetVerificationId(email string) (string, error) {
	var verificationId string
	err := db.Model(&ResetPassword{}).Where("email = ?", email).Select("verification_id").First(&verificationId).Error
	if err != nil {
		return "", err
	}
	return verificationId, nil
}

func VerifyVerificationId(email, verificationId string) (bool, error) {
	var Correct int64
	err := db.Model(&ResetPassword{}).Where("email = ? AND verification_id = ?", email, verificationId).Count(&Correct).Error
	if err != nil {
		return false, err
	}

	return Correct >= 1, nil
}

func CreateResetPassword(email string) error {
	guid := xid.New()
	verificationId := guid.String()

	resetPwdObj := ResetPassword{
		VerificationID:  verificationId,
		Email:           email,
		AlreadyVerified: false,
	}

	err := db.Create(&resetPwdObj).Error
	if err != nil {
		return err
	}

	err = sendEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(email string) error {
	//TODO: SEND EMAIL
	return nil
}
