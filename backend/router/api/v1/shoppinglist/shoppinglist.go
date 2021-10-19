package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
	"github.com/urento/shoppinglist/services"
)

func GetShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

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
		})
		return
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"success": "false",
		})
		return
	}

	listService := services.Shoppinglist{ID: id, Owner: owner}
	list, err := listService.GetList()
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_LIST_FAIL, map[string]string{
			"success": "false",
		})
		return
	}

	if list.Owner != owner {
		appG.Response(http.StatusBadRequest, e.ERROR_LIST_DOES_NOT_BELONG_TO_TOKEN, map[string]string{
			"error":   "list does not belong to reuqest maker",
			"success": "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

func GetShoppinglists(c *gin.Context) {
	appG := app.Gin{C: c}

	var o int
	offset, ok := c.GetQuery("offset")
	if !ok {
		o = 0
	} else {
		offsetToInt, err := strconv.Atoi(offset)
		if err != nil {
			log.Print(err)
			appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
				"error":   err.Error(),
				"success": "false",
			})
			return
		}
		o = offsetToInt
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

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	list := services.Shoppinglist{Owner: email}
	lists, err := list.GetListsByOwner(o)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_LISTS_BY_OWNER, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

func GetShoppinglistsByParticipation(c *gin.Context) {
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

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	lists, err := services.GetListsByParticipant(email)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_LISTS_BY_OWNER, map[string]string{
			"success": "false",
			"error":   "error while getting lists by participation",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

type CreateShoppinglistForm struct {
	Title string `form:"title"`
}

func CreateShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	var f CreateShoppinglistForm

	if err := c.BindJSON(&f); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	//TODO: Validate data some other way

	/*if f.Participants == nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}*/

	if f.Title == "" {
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

	randomId := util.RandomIntWithLength(900000)
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

	lists := services.Shoppinglist{
		ID:    randomId,
		Title: f.Title,
		Owner: owner,
		//Participants: f.Participants,
	}

	if _, err := lists.Create(userId, true); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_LIST_FAIL, map[string]string{
			"success": "false",
			"message": "Error while creating Shoppinglist",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

type EditShoppinglistForm struct {
	ID           int                   `form:"id"`
	Title        string                `form:"title" valid:"Required"`
	Items        models.Item           `form:"items"`
	Owner        string                `form:"owner" valid:"Required"`
	Participants []*models.Participant `form:"participants" valid:"Required"`
}

func EditShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	var form EditShoppinglistForm

	if err := c.BindJSON(&form); err != nil {
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

	//TODO: Validate data some other way

	if form.Owner == "" {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	if form.Title == "" {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	list := services.Shoppinglist{
		ID:           form.ID,
		Title:        form.Title,
		Items:        form.Items,
		Owner:        form.Owner,
		Participants: form.Participants,
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

	err = list.Edit(userId, false)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_LIST_FAIL, map[string]string{"success": "false"})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

func DeleteShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, map[string]string{
			"success": "false",
			"message": "validation error",
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

	listService := services.Shoppinglist{ID: id}
	exists, err := listService.ExistsByID()
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_LIST_FAIL, map[string]string{
			"success": "false",
			"message": "list does not exist",
		})
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_LIST_DOES_NOT_EXIST, map[string]string{
			"success": "false",
			"message": "list does not exist",
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

	err = listService.Delete(userId, true)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_LIST_FAIL, map[string]string{
			"success": "false",
			"message": "error deleting the list",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

type ItemRequest struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Position int    `json:"position"`
}

func AddItem(c *gin.Context) {
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

	var form ItemRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	if form.Title == "" {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "title can't be empty",
			"success": "false",
		})
		return
	}

	owner, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
		})
		return
	}

	itemId := util.RandomIntWithLength(900000)
	id := form.ID
	item := &models.Item{
		ParentListID: id,
		ItemID:       itemId,
		Title:        form.Title,
		Position:     int64(form.Position),
		Bought:       false,
	}
	shoppinglist := &services.Shoppinglist{
		ID:    id,
		Items: *item,
		Owner: owner,
	}

	item, err = shoppinglist.AddItem()
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while adding item",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, item)
}

type UpdateItemRequest struct {
	ParentListID int    `json:"parentListId"`
	Title        string `json:"title"`
	Position     int64  `json:"position"`
	Bought       bool   `json:"bought"`
}

func UpdateItem(c *gin.Context) {
	appG := app.Gin{C: c}
	itemId := com.StrTo(c.Param("id")).MustInt()
	var form UpdateItemRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	item := models.Item{
		ParentListID: form.ParentListID,
		ItemID:       itemId,
		Title:        form.Title,
		Position:     form.Position,
		Bought:       form.Bought,
	}

	err := models.UpdateItem(item)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while updating item",
			"success": "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, item)
}

type UpdateItemsRequest struct {
	ParentListID int           `json:"parent_list_id"`
	Items        []models.Item `json:"items"`
}

func UpdateItems(c *gin.Context) {
	appG := app.Gin{C: c}
	var form UpdateItemsRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	err := models.UpdateItems(form.ParentListID, form.Items)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while updating items",
			"success": "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}

type DeleteItemRequest struct {
	ID           int `json:"id"`
	ParentListId int `json:"parent_list_id"`
}

func DeleteItem(c *gin.Context) {
	appG := app.Gin{C: c}
	var form DeleteItemRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   "error while binding json to struct",
			"success": "false",
		})
		return
	}

	err := models.DeleteItem(form.ParentListId, form.ID)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while deleting item",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})

}

func GetListItems(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id")

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
		})
		return
	}

	shoppinglist := services.Shoppinglist{
		ID:    id,
		Owner: email,
	}

	items, err := shoppinglist.GetItems()
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_LISTS_BY_OWNER, map[string]string{
			"success": "false",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, items)
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
