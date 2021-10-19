package models

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/util"
)

type BackupCodes struct {
	Model

	Owner string         `json:"owner"`
	Codes pq.StringArray `gorm:"type:text[]" json:"codes"`
}

func GenerateCodes(email string, userId int, regenerate bool, withNotification bool) (pq.StringArray, error) {
	code1 := util.RandomString(8)
	code2 := util.RandomString(8)
	code3 := util.RandomString(8)
	code4 := util.RandomString(8)
	code5 := util.RandomString(8)
	code6 := util.RandomString(8)
	codes := pq.StringArray{
		code1,
		code2,
		code3,
		code4,
		code5,
		code6,
	}

	backupCodes := &BackupCodes{
		Owner: email,
		Codes: codes,
	}

	if regenerate {
		has, err := HasCodes(email)
		if err != nil {
			return pq.StringArray{}, err
		}

		if has {
			err := RemoveCodes(email)
			if err != nil {
				return pq.StringArray{}, err
			}
		}
	}

	if withNotification {
		notification := Notification{
			UserID:           userId,
			Title:            "New Backup Codes",
			Text:             "New Backup Codes were generated",
			NotificationType: "new_backupcodes",
			Date:             time.Now().Format("02.01.2006"),
		}

		if err := CreateNotification(notification); err != nil {
			return pq.StringArray{}, err
		}
	}

	err := db.Create(&backupCodes).Error
	if err != nil {
		return pq.StringArray{}, err
	}

	return codes, nil
}

func GetCodes(email string) (pq.StringArray, error) {
	var Codes pq.StringArray
	err := db.Model(&BackupCodes{}).Where("owner = ?", email).Select("codes").Find(&Codes).Error
	return Codes, err
}

func RemoveCodes(email string) error {
	err := db.Where("owner = ?", email).Delete(&BackupCodes{Owner: email}).Error
	return err
}

func HasCodes(email string) (bool, error) {
	var Has bool
	err := db.Raw("SELECT EXISTS(SELECT codes FROM backup_codes WHERE owner = ?) AS found", email).Scan(&Has).Error
	return Has, err
}

func VerifyCode(email, code string) (bool, error) {
	var Codes []string
	err := db.Model(&BackupCodes{}).Select("codes").Where("owner = ?", email).Find(&Codes).Error
	if err != nil {
		return false, err
	}

	for idx := range Codes {
		b := util.StringArrayToArray(Codes, idx)
		if b == code {
			err = cache.ActivateResetPassword(context.Background(), email)
			if err != nil {
				return true, err
			}
			return true, nil
		}
	}

	return false, nil
}
