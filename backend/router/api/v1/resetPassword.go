package api

import (
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
	"github.com/urento/shoppinglist/services"
)

type ResetPassword struct {
	Email string `valid:"Required"`
}

type VerifyId struct {
	Email          string `valid:"Required"`
	VerificationId string `valid:"Required"`
}

func SendResetPassword(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	var resetPassword ResetPassword

	if err := c.BindJSON(&resetPassword); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email := resetPassword.Email

	emailByJwt, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	if emailByJwt != email {
		appGin.Response(http.StatusBadRequest, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"message": "wrong email",
			"success": "false",
		})
		return
	}

	ok, err := valid.Valid(&ResetPassword{Email: email})
	if !ok {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	rpwd := services.ResetPassword{Email: email}
	err = rpwd.CreateResetPassword()
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_SENDING_RESET_PASSWORD_EMAIL, map[string]string{
			"message": "error while sending reset password email",
			"success": "false",
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"message": "successfully sent the reset password email",
		"success": "true",
	})
}

func VerifyVerificationId(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	var resetPassword VerifyId

	if err := c.BindJSON(&resetPassword); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email := resetPassword.Email
	verificationId := resetPassword.VerificationId

	emailByJwt, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	if emailByJwt != email {
		appGin.Response(http.StatusBadRequest, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"message": "Wrong Email",
			"success": "false",
		})
		return
	}

	resetPwdObj := VerifyId{
		Email:          email,
		VerificationId: verificationId,
	}

	ok, err := valid.Valid(&resetPwdObj)
	if !ok {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	rpwd := services.ResetPassword{
		Email:          email,
		VerificationId: verificationId,
	}

	correct, err := rpwd.VerifyVerificationId()
	if err != nil || !correct {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.ERROR_VERIFYING_VERIFICATION_ID, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"message": "Verification ID is correct",
	})
}

func ChangePassword(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	var resetPassword VerifyId

	if err := c.BindJSON(&resetPassword); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email := resetPassword.Email
	verificationId := resetPassword.VerificationId

	emailByJwt, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	if emailByJwt != email {
		appGin.Response(http.StatusBadRequest, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"message": "Wrong Email",
			"success": "false",
		})
		return
	}

	resetPwdObj := VerifyId{
		Email:          email,
		VerificationId: verificationId,
	}

	ok, err := valid.Valid(&resetPwdObj)
	if !ok {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//TODO: Finish this

}
