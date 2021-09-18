package services

import (
	"github.com/urento/shoppinglist/models"
)

type Auth struct {
	EMail                   string
	Username                string
	Password                string
	EmailVerified           bool
	Rank                    string
	JWTToken                string
	TwoFactorAuthentication bool
	IPAddress               string
	Disabled                bool
}

func (auth *Auth) Check() (bool, error) {
	return models.CheckAuth(auth.EMail, auth.Password, auth.IPAddress)
}

func (auth *Auth) Create() error {
	return models.CreateAccount(auth.EMail, auth.Username, auth.Password, auth.IPAddress)
}

func (auth *Auth) UpdateIP() error {
	return models.UpdateIP(auth.EMail, auth.IPAddress)
}

func (auth *Auth) GetIP() (string, error) {
	return models.GetIP(auth.EMail)
}

func (auth *Auth) VerifyEmail() error {
	return models.VerifyEmail(auth.EMail)
}

func (auth *Auth) IsEmailVerified() (bool, error) {
	return models.IsEmailVerified(auth.EMail)
}

func (auth *Auth) SendVerificationEmail() error {
	return models.SendVerifyEmail(auth.EMail)
}

func (auth *Auth) SetRank() error {
	return models.SetRank(auth.EMail, auth.Rank)
}

func (auth *Auth) GetRank() (string, error) {
	return models.GetRank(auth.EMail)
}

func (auth *Auth) Delete() error {
	return models.DeleteAccount(auth.EMail, auth.Password)
}

func (auth *Auth) GetUser() (*models.Auth, error) {
	return models.GetUser(auth.EMail)
}

func (auth *Auth) Logout() (bool, error) {
	return models.Logout(auth.EMail, auth.JWTToken)
}

func (auth *Auth) GetPassword() (string, error) {
	return models.GetPasswordHash(auth.EMail)
}

func (auth *Auth) SetTwoFactorAuthentication() error {
	return models.SetTwoFactorAuthentication(auth.EMail, auth.TwoFactorAuthentication)
}

func (auth *Auth) IsTwoFactorEnabled() (bool, error) {
	return models.IsTwoFactorEnabled(auth.EMail)
}

func (auth *Auth) SetUsername() error {
	return models.SetUsername(auth.EMail, auth.Username)
}

func (auth *Auth) GetUsername() (string, error) {
	return models.GetUsername(auth.EMail)
}

func (auth *Auth) DisableAccount() error {
	return models.DisableAccount(auth.EMail)
}

func (auth *Auth) ActivateAccount() error {
	return models.ActivateAccount(auth.EMail)
}

func (auth *Auth) IsDisabled() (bool, error) {
	return models.IsDisabled(auth.EMail)
}
