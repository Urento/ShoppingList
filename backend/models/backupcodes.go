package models

import (
	"math/rand"
	"time"

	"github.com/lib/pq"
)

type BackupCodes struct {
	Model

	Owner string         `json:"owner"`
	Codes pq.StringArray `gorm:"type:text[]" json:"codes"`
}

func GenerateCodes(email string, regenerate bool) (pq.StringArray, error) {
	code1 := GenerateCode()
	code2 := GenerateCode()
	code3 := GenerateCode()
	code4 := GenerateCode()
	code5 := GenerateCode()
	codes := pq.StringArray{
		code1,
		code2,
		code3,
		code4,
		code5,
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

	err := db.Create(&backupCodes).Error
	if err != nil {
		return pq.StringArray{}, err
	}

	return codes, nil
}

//TODO: ADD TEST
func GetCodes(email string) (pq.StringArray, error) {
	var Codes pq.StringArray
	err := db.Debug().Model(&BackupCodes{}).Where("owner = ?", email).Select("codes").Find(&Codes).Error
	return Codes, err
}

func RemoveCodes(email string) error {
	err := db.Debug().Where("owner = ?", email).Delete(&BackupCodes{Owner: email}).Error
	return err
}

func HasCodes(email string) (bool, error) {
	var Has bool
	err := db.Raw("SELECT EXISTS(SELECT codes FROM backup_codes WHERE owner = ?) AS found", email).Scan(&Has).Error
	return Has, err
}

func VerifyCode(email, code string) (bool, []string, error) {
	var codes []string
	err := db.Debug().Model(&BackupCodes{}).Where("owner = ?", email).Select("codes").Scan(&codes).Error

	return false, codes, err
}

func GenerateCode() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
