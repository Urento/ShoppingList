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

func GetNotifications(c *gin.Context) {
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

	notifications, err := models.GetNotifications(userId)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting notifications",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, notifications)
}

type NotificationRequest struct {
	NotificationId int `json:"notification_id"`
}

func DeleteNotification(c *gin.Context) {
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

	var f NotificationRequest

	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	err = models.DeleteNotification(userId, f.NotificationId)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while deleting notification",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}

func GetNotification(c *gin.Context) {
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

	notificationId := c.Param("notification_id")
	notificationIdAsInt, err := strconv.Atoi(notificationId)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	notification, err := models.GetNotification(userId, notificationIdAsInt)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"success": "false",
			"error":   "error while getting notification",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, notification)
}
