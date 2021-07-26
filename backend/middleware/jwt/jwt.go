package jwt

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int

		code = e.SUCCESS
		token := c.Request.Header.Get("Authorization")
		splitToken := strings.Replace(token, "Bearer ", "", -1)
		if splitToken == "" {
			code = e.ERROR_NOT_AUTHORIZED
		} else {
			//check if token is valid in redis
			tokenValid, err := cache.IsTokenValid(splitToken)
			if err != nil || !tokenValid {
				log.Print(err)
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}

			if tokenValid {
				_, err := util.ParseToken(splitToken)

				if err != nil {
					switch err.(*jwt.ValidationError).Errors {
					case jwt.ValidationErrorExpired:
						code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					default:
						code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
					}
				}
			}
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
