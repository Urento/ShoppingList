package services

import "github.com/urento/shoppinglist/models"

type Auth struct {
	EMail         string
	Username      string
	Password      string
	EmailVerified bool
	Rank          string
	JWTToken      string
}

func (auth *Auth) Check() (bool, error) {
	return models.CheckAuth(auth.EMail, auth.Password)
}

func (auth *Auth) Create() error {
	return models.CreateAccount(auth.EMail, auth.Username, auth.Password)
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
