package v1

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/services"
)

func GetShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	//TODO: Check if the Shoppinglist belongs to the request maker
	tok := c.Request.Header.Get("Authentication")
	token := strings.Replace(tok, "Bearer ", "", -1)

	emailOfRequestMaker, err := cache.GetEmailByJWT(token)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	listService := services.Shoppinglist{ID: id}
	exists, err := listService.ExistsByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_LIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_LIST_DOES_NOT_EXIST, nil)
		return
	}

	list, err := listService.GetList()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_LIST_FAIL, nil)
		return
	}

	if list.Owner != emailOfRequestMaker {
		appG.Response(http.StatusBadRequest, e.ERROR_LIST_DOES_NOT_BELONG_TO_TOKEN, map[string]string{
			"error": "list does not belong to reuqest maker",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

func GetShoppinglists(c *gin.Context) {
	appG := app.Gin{C: c}
	tok := c.Request.Header.Get("Authorization")
	token := strings.Replace(tok, "Bearer ", "", -1)

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	listService := services.Shoppinglist{Owner: email}
	lists, err := listService.GetListsByOwner()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_LISTS_BY_OWNER, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

type CreateShoppinglistForm struct {
	ID           int      `form:"id"`
	Title        string   `form:"title" valid:"Required"`
	Items        []string `form:"items" valid:"Required"`
	Owner        string   `form:"owner" valid:"Required"`
	Position     int      `form:"position" valid:"Required"`
	Participants []string `form:"participants" valid:"Required"`
}

func CreateShoppinglist(c *gin.Context) {
	var (
		appG       = app.Gin{C: c}
		form       CreateShoppinglistForm
		seededRand *rand.Rand = rand.New(
			rand.NewSource(time.Now().UnixNano()))
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, map[string]string{
			"error": "validation error",
		})
		return
	}

	randomId := seededRand.Intn(900000)

	lists := services.Shoppinglist{
		ID:           randomId,
		Title:        form.Title,
		Items:        form.Items,
		Owner:        form.Owner,
		Participants: form.Participants,
		Position:     form.Position,
	}

	exists, err := lists.ExistsByID()
	if err != nil || exists {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_LIST_DOES_NOT_EXIST, map[string]string{
			"id":    strconv.Itoa(randomId),
			"error": "id already in use - retry",
		})
		return
	}

	if _, err := lists.Create(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

type EditShoppinglistForm struct {
	ID           int      `form:"id"`
	Title        string   `form:"title" valid:"Required"`
	Items        []string `form:"items" valid:"Required"`
	Owner        string   `form:"owner" valid:"Required"`
	Position     int      `form:"position" valid:"Required"`
	Participants []string `form:"participants" valid:"Required"`
}

func EditShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_ID_IS_INVALID, map[string]string{
			"error": "id is invalid",
		})
		return
	}

	var (
		form = EditShoppinglistForm{ID: id}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, map[string]string{
			"error": "validation error",
		})
		return
	}

	//TODO: Validate all Emails in Participants
	list := services.Shoppinglist{
		ID:           id,
		Title:        form.Title,
		Items:        form.Items,
		Owner:        form.Owner,
		Participants: form.Participants,
		Position:     form.Position,
	}
	exists, err := list.ExistsByID()
	if err != nil || !exists {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_CHECK_EXIST_LIST_FAIL, nil)
		return
	}

	err = list.Edit()
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

func DeleteShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id")
	log.Println(id)

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, map[string]string{
			"success": "false",
			"message": "validation error",
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

	err = listService.Delete()
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
