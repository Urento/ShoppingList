package notifications_v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
)

func HasUnreadNotifications(c *gin.Context) {
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

	has, err := models.HasUnreadNotifications(userId)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECKING_HAS_UNREAD_NOTIFICATIONS, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"has":     strconv.FormatBool(has),
	})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
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

	err = models.MarkAllNotificationsAsRead(userId)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while marking all notifications as read",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}
