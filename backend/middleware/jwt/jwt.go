package jwt

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int

		code = e.SUCCESS
		token, err := GetCookie(c)
		if err != nil {
			code = e.ERROR_GETTING_HTTPONLY_COOKIE
		}

		if token == "" {
			code = e.ERROR_NOT_AUTHORIZED
		} else {
			//check if token is valid in redis
			tokenValid, err := cache.IsTokenValid(token)
			if err != nil || !tokenValid {
				log.Print(err)
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}

			if tokenValid {
				data, parseErr := util.ParseToken(token)

				ok, err := cache.VerifySecretId(data.Email, data.SecretId)
				if err != nil || !ok {
					log.Print(err)
					code = e.ERROR_VERIFYING_VERIFICATION_ID
				}

				if parseErr != nil {
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					default:
						code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
					}
				}
			}
		}

		if code == e.ERROR_VERIFYING_VERIFICATION_ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    "jwt token is invalid",
			})
			c.Abort()
			return
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    map[string]string{"token": token},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetCookie(ctx *gin.Context) (string, error) {
	token, err := ctx.Request.Cookie("token")
	if err != nil {
		return "", err
	}

	if len(token.Value) <= 0 {
		return "", errors.New("cookie 'token' has to be longer than 0 charcters")
	}

	if len(token.Value) <= 50 {
		return "", errors.New("cookie 'token' has to be longer than 50 charcters")
	}

	return token.Value, nil
}
