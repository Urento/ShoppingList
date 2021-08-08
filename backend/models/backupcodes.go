package models

import (
	"math/rand"
	"time"
)

type BackupCodes struct {
	Model

	Owner  string `json:"owner"`
	Code1  string `json:"code1"`
	Code2  string `json:"code2"`
	Code3  string `json:"code3"`
	Code4  string `json:"code4"`
	Code5  string `json:"code5"`
	Code6  string `json:"code6"`
	Code7  string `json:"code7"`
	Code8  string `json:"code8"`
	Code9  string `json:"code9"`
	Code10 string `json:"code10"`
}

type Codes struct {
	Code1  string `json:"code1"`
	Code2  string `json:"code2"`
	Code3  string `json:"code3"`
	Code4  string `json:"code4"`
	Code5  string `json:"code5"`
	Code6  string `json:"code6"`
	Code7  string `json:"code7"`
	Code8  string `json:"code8"`
	Code9  string `json:"code9"`
	Code10 string `json:"code10"`
}

func GenerateCodes(email string) (*Codes, error) {
	code1 := GenerateCode()
	code2 := GenerateCode()
	code3 := GenerateCode()
	code4 := GenerateCode()
	code5 := GenerateCode()
	code6 := GenerateCode()
	code7 := GenerateCode()
	code8 := GenerateCode()
	code9 := GenerateCode()
	code10 := GenerateCode()
	backupCodes := BackupCodes{
		Owner:  email,
		Code1:  code1,
		Code2:  code2,
		Code3:  code3,
		Code4:  code4,
		Code5:  code5,
		Code6:  code6,
		Code7:  code7,
		Code8:  code8,
		Code9:  code9,
		Code10: code10,
	}

	hasCodes, err := HasCodes(email)
	if err != nil {
		return &Codes{}, err
	}

	if hasCodes {
		err := RemoveOldCodes(email)
		if err != nil {
			return &Codes{}, err
		}
	}

	if err := db.Model(&BackupCodes{}).Create(&backupCodes).Error; err != nil {
		return &Codes{}, err
	}

	return &Codes{Code1: code1, Code2: code2, Code3: code3, Code4: code4, Code5: code5, Code6: code6, Code7: code7, Code8: code8, Code9: code9, Code10: code10}, nil
}

func GetCodes(email string) (*Codes, error) {
	var codes Codes
	err := db.Debug().Model(&BackupCodes{}).Where("owner = ?", email).Select([]string{"code1", "code2", "code3", "code4", "code5", "code6", "code7", "code8", "code9", "code10"}).Find(&codes).Error
	return &codes, err
}

func RemoveCode(email, code string) error {
	return nil
}

func RemoveOldCodes(email string) error {
	return nil
}

func HasCodes(email string) (bool, error) {
	var has bool
	err := db.Raw("SELECT EXISTS(SELECT code1 FROM backup_codes WHERE owner = ?) AS found", email).Scan(&has).Error
	return has, err
}

func VerifyCode(email, code string) (bool, error) {
	return false, nil
}

func GenerateCode() string {
	length := 8
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
