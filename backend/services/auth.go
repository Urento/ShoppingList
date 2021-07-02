package services

import "github.com/urento/shoppinglist/models"

type Auth struct {
	Username string
	Password string
}

func (auth *Auth) Check() (bool, error) {
	return models.CheckAuth(auth.Username, auth.Password)
}
