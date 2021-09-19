package services

import "github.com/urento/shoppinglist/models"

type ResetPassword struct {
	Email          string `json:"email"`
	VerificationId string `json:"verification_id"`
}

func (rpwd *ResetPassword) ExistsResetPassword() (bool, error) {
	return models.HasResetPassword(rpwd.Email)
}

func (rpwd *ResetPassword) DeleteResetPassword() error {
	return models.DeleteResetPassword(rpwd.Email)
}

func (rpwd *ResetPassword) CreateResetPassword() error {
	return models.CreateResetPassword(rpwd.Email)
}

func (rpwd *ResetPassword) VerifyVerificationId() (bool, error) {
	return models.VerifyVerificationID(rpwd.Email, rpwd.VerificationId)
}
