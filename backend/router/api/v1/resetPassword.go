package api

import (
	"context"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
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

	err = models.CreateResetPassword(email)
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

	correct, err := models.VerifyVerificationID(email, verificationId)
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

type ChangePasswordRequest struct {
	Code     string `json:"code"`
	Owner    string `json:"owner"`
	Password string `json:"password"`
}

func ChangePassword(c *gin.Context) {
	appG := app.Gin{C: c}
	var form ChangePasswordRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, map[string]string{
			"success":  "false",
			"verified": "false",
			"ok":       "false",
		})
		return
	}

	ok, err := models.VerifyCode(form.Owner, form.Code)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while verifying backupcode",
			"success": "false",
			"ok":      "false",
		})
		return
	}

	if !ok {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"error":   "code is incorrect",
			"success": "false",
			"ok":      "false",
		})
		return
	}

	canReset, err := cache.CanResetPassword(context.Background(), form.Owner)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while checking if the user can reset password",
			"ok":      "false",
		})
		return
	}

	if !canReset {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "you can't reset your password",
			"ok":      "false",
		})
		return
	}

	err = models.ResetPasswordFromUser(form.Owner, form.Password, "", false)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while resetting password from user without old password",
			"ok":      "true",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"ok":      "true",
	})
}
