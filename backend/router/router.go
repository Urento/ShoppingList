package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/middleware/jwt"
	"github.com/urento/shoppinglist/middleware/ratelimiter"
	"github.com/urento/shoppinglist/router/api/v1"
	v1 "github.com/urento/shoppinglist/router/api/v1/shoppinglist"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(ratelimiter.Ratelimiter())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"OPTIONS", "PUT", "GET", "POST", "DELETE", "PATCH"},
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           9 * time.Hour,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	r.POST("/api/auth", api.Login)
	r.POST("/api/auth/register", api.CreateAccount)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())

	apiv1.POST("/auth/check", api.Check)
	apiv1.GET("/auth/user", api.GetUser)
	apiv1.POST("/auth/logout", api.Logout)
	apiv1.POST("/auth/update", api.UpdateUser)

	apiv1.GET("/lists", v1.GetShoppinglists)
	apiv1.POST("/list", v1.CreateShoppinglist)
	apiv1.PUT("/list", v1.EditShoppinglist)
	apiv1.GET("/list/:id", v1.GetShoppinglist)
	apiv1.DELETE("/list/:id", v1.DeleteShoppinglist)

	apiv1.POST("/resetpassword/verifyid", api.VerifyVerificationId)
	apiv1.POST("/resetpassword", api.SendResetPassword)
	apiv1.POST("/resetpassword/changepassword", api.ChangePassword)

	apiv1.POST("/twofactorauthentication", api.UpdateTwoFactorAuthentication)

	return r
}
