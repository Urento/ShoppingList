package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
	auth "github.com/urento/shoppinglist/services"
)

type Auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required"`
}

func GetAuth(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := Auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}

	authService := auth.Auth{Username: username, Password: password}
	exists, err := authService.Check()
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !exists {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
