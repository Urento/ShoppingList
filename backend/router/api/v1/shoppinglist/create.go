package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/e"
)

type ShoppinglistForm struct {
	Id           string   `form:"id"`
	Title        string   `form:"title" valid:"Required"`
	Items        []string `form:"items" valid:"Required"`
	Owner        string   `form:"owner" valid:"Required"`
	Participants []string `form:"participants" valid:"Required"`
}

func CreateShoppinglist(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form ShoppinglistForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
}
