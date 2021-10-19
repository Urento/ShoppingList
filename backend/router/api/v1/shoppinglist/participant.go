package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
)

type AddParticipantRequest struct {
	Email        string `json:"email"`
	ParentListId int    `json:"parentListId"`
}

func AddParticipant(c *gin.Context) {
	appG := app.Gin{C: c}
	var f AddParticipantRequest

	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	included, err := models.IsParticipantAlreadyIncluded(f.Email, f.ParentListId)
	if err != nil || included {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "participant is already included",
			"success": "false",
		})
		return
	}

	p := models.Participant{
		ParentListID: f.ParentListId,
		Email:        f.Email,
		Status:       "pending",
	}

	participant, err := models.AddParticipant(p)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while adding participant",
			"success": "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, participant)
}

func GetPendingRequests(c *gin.Context) {
	appG := app.Gin{C: c}

	email := c.Param("email")

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "",
		})
		return
	}

	if owner != email {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "ja ja ja",
		})
		return
	}

	requests, err := models.GetPendingRequests(email)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting pending requests",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, requests)
}

type AcceptRequestRequest struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func AcceptRequest(c *gin.Context) {
	appG := app.Gin{C: c}
	var f AcceptRequestRequest

	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting email by jwt",
		})
		return
	}

	if owner != f.Email {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "ja ja ja",
		})
		return
	}

	err = models.AcceptRequest(f.ID, f.Email)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while accepting the request",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

type DeleteRequestRequest struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func DeleteRequest(c *gin.Context) {
	appG := app.Gin{C: c}
	var f DeleteRequestRequest

	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "",
		})
		return
	}

	if owner != f.Email {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "ja ja ja",
		})
		return
	}

	err = models.DeleteRequest(f.ID, f.Email)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while deleting the request",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

func GetGetPendingRequestsFromShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "",
		})
		return
	}

	requests, err := models.GetPendingRequestsFromShoppinglist(owner, id)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting pending requests",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, requests)
}

func GetParticipants(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting email by jwt",
		})
		return
	}

	belongs, err := models.BelongsShoppinglistToEmail(owner, id)
	if err != nil || !belongs {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "shoppinglist doesn't belong to the given jwt token",
		})
		return
	}

	participants, err := models.GetParticipants(id)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting participants",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, participants)
}

func DeleteParticipant(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	parentListId := com.StrTo(c.Param("parentListId")).MustInt()

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
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while getting email by jwt",
		})
		return
	}

	belongs, err := models.BelongsShoppinglistToEmail(owner, parentListId)
	if err != nil || !belongs {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "shoppinglist doesn't belong to the given jwt token",
		})
		return
	}

	err = models.RemoveParticipant(parentListId, id)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while deleting participant",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}
