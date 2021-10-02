package util

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/urento/shoppinglist/pkg/cache"
)

var jwtSecret []byte

type Claims struct {
	Email    string `json:"email"`
	SecretId string `json:"secretId"`
	jwt.StandardClaims
}

func GenerateToken(email string, refreshToken bool) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	refreshTokenExpireTime := nowTime.Add(168 * time.Hour) // 1 week in hours

	secretId, err := cache.GenerateSecretId(email)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		email,
		secretId,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "shoppinglist",
		},
	}

	if refreshToken {
		claims = &Claims{
			email,
			secretId,
			jwt.StandardClaims{
				ExpiresAt: refreshTokenExpireTime.Unix(),
				Issuer:    "shoppinglist",
			},
		}
	} else {
		claims = &Claims{
			email,
			secretId,
			jwt.StandardClaims{
				ExpiresAt: expireTime.Unix(),
				Issuer:    "shoppinglist",
			},
		}
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
