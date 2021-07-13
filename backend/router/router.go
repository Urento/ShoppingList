package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/middleware/cors"
	"github.com/urento/shoppinglist/middleware/jwt"
	"github.com/urento/shoppinglist/router/api/v1"
	v1 "github.com/urento/shoppinglist/router/api/v1/shoppinglist"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/api/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(cors.CORSMiddleware())
	apiv1.Use(jwt.JWT())
	apiv1.GET("/lists", v1.GetShoppinglists)
	apiv1.POST("/list", v1.CreateShoppinglist)
	apiv1.GET("/list/:id", v1.GetShoppinglist)
	apiv1.DELETE("/list/:id", v1.DeleteShoppinglist)

	return r
}
