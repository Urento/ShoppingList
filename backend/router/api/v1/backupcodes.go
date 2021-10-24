package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
)

type VerifyBackupCodes struct {
	Owner string `json:"owner"`
	Code  string `json:"code"`
}

func VerifyBackupCode(c *gin.Context) {
	appG := app.Gin{C: c}

	var f VerifyBackupCodes
	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
			"ok":      "false",
		})
		return
	}

	//TODO: add custom ratelimiter and/or some other form of security

	ok, err := models.VerifyCode(f.Owner, f.Code)
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

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"ok":      "true",
	})
}

func GenerateCodes(c *gin.Context) {
	appG := app.Gin{C: c}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	owner, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	userId, err := models.GetUserIDByEmail(owner)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	codes, err := models.GenerateCodes(owner, userId, false, true)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while regenerating codes",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"codes":   strings.Join(codes, ","),
	})
}

func RegenerateCodes(c *gin.Context) {
	appG := app.Gin{C: c}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	owner, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	userId, err := models.GetUserIDByEmail(owner)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	codes, err := models.GenerateCodes(owner, userId, true, true)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while regenerating codes",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"codes":   strings.Join(codes, ","),
	})
}

func GetBackupCodes(c *gin.Context) {
	appG := app.Gin{C: c}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	owner, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	has, err := models.HasCodes(owner)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while checking if the owner has backupcodes",
			"has":     "false",
		})
		return
	}

	if !has {
		log.Print(err)
		appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success": "false",
			"has":     "false",
		})
		return
	}

	codes, err := models.GetCodes(owner)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while getting codes",
			"has":     "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"codes":   strings.Join(codes, ","),
		"has":     "true",
	})
}
