package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
	auth "github.com/urento/shoppinglist/services"
)

type Auth struct {
	Email    string `valid:"Required;"`
	Username string
	Password string `valid:"Required"`
}

func Check(c *gin.Context) {
	appGin := app.Gin{C: c}
	token := c.Request.Header.Get("Authorization")

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"success": "false",
		})
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
		})
	}

	check, err := cache.Check(email, token)
	if !check || err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_TOKEN_INVALID, map[string]string{
			"success": "false",
		})
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

func GetUser(c *gin.Context) {
	appGin := app.Gin{C: c}
	token := c.Request.Header.Get("Authorization")

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
	}

	authService := auth.Auth{EMail: email}
	data, err := authService.GetUser()
	if err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_RETRIEVING_USER_DATA, nil)
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

func GetAuth(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	email := c.PostForm("email")
	password := c.PostForm("password")

	a := Auth{Email: email, Password: password}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	authService := auth.Auth{EMail: email, Password: password}
	exists, err := authService.Check()
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !exists {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(email, password)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

func CreateAccount(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	a := Auth{Email: email, Username: username, Password: password}
	ok, _ := valid.Valid(a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "validation error",
		})
		return
	}

	acc := auth.Auth{EMail: email, Username: username, Password: password, EmailVerified: false, Rank: "default"}
	err := acc.Create()
	if err != nil {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.ERROR_CREATING_ACCOUNT, err)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"created":  "true",
		"email":    email,
		"username": username,
	})
}
