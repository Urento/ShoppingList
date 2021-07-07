package services

import "github.com/urento/shoppinglist/models"

type Auth struct {
	EMail    string
	Username string
	Password string
}

func (auth *Auth) Check() (bool, error) {
	if auth.Username != "" {
		return models.CheckAuth("", auth.Username, auth.Password)
	}
	return models.CheckAuth(auth.EMail, "", auth.Password)
}

func (auth *Auth) Create() error {
	return models.CreateAccount(auth.EMail, auth.Username, auth.Password)
}
