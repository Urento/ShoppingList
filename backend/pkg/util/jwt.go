package util

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/dgrijalva/jwt-go"
	"github.com/urento/shoppinglist/pkg/cache"
)

var jwtSecret []byte

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(email, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	claims := Claims{
		email,
		passwordHash,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "shoppinglist",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	err = cache.CacheJWT(email, token)
	if err != nil {
		return "", err
	}

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
