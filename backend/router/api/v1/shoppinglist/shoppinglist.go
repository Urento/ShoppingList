package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/services"
)

func GetShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

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
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_LIST, nil)
		return
	}

	list, err := listService.GetList()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, list)
}

//TODO
func GetShoppinglists(c *gin.Context) {

}

type CreateAndEditShoppinglistForm struct {
	ID           int      `form:"id"`
	Title        string   `form:"title" valid:"Required"`
	Items        []string `form:"items" valid:"Required"`
	Owner        string   `form:"owner" valid:"Required"`
	Participants []string `form:"participants" valid:"Required"`
}

func CreateShoppinglist(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form CreateAndEditShoppinglistForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	listService := services.Shoppinglist{ID: form.ID}
	exists, err := listService.ExistsByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_LIST, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_LIST, nil)
		return
	}

	lists := services.Shoppinglist{
		ID:           form.ID,
		Title:        form.Title,
		Items:        form.Items,
		Owner:        form.Owner,
		Participants: form.Participants,
	}
	if err := lists.Create(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, lists)
}

func EditShoppinglist(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = CreateAndEditShoppinglistForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	listService := services.Shoppinglist{
		ID:           form.ID,
		Title:        form.Title,
		Items:        form.Items,
		Owner:        form.Owner,
		Participants: form.Participants,
	}
	exists, err := listService.ExistsByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_LIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_LIST, nil)
		return
	}

	err = listService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func DeleteShoppinglist(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	listService := services.Shoppinglist{ID: id}
	exists, err := listService.ExistsByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_LIST_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_LIST, nil)
		return
	}

	err = listService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_LIST_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
