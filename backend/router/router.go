package routers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/middleware/jwt"
	"github.com/urento/shoppinglist/middleware/ratelimiter"
	"github.com/urento/shoppinglist/router/api/v1"
	notifications_v1 "github.com/urento/shoppinglist/router/api/v1/notifications"
	v1 "github.com/urento/shoppinglist/router/api/v1/shoppinglist"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	// maybe remove because gin attaches them automatically
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
	//apiv1.POST("/auth/invalidate", api.InvalidateSpecificJWTToken) //TODO: Test this and add this to the frontend

	apiv1.GET("/lists", v1.GetShoppinglists)
	apiv1.GET("/listsByParticipation", v1.GetShoppinglistsByParticipation)
	apiv1.POST("/list", v1.CreateShoppinglist)
	apiv1.PUT("/list", v1.EditShoppinglist)
	apiv1.GET("/list/:id", v1.GetShoppinglist)
	apiv1.GET("/list/items/:id", v1.GetListItems) //TODO: Start using this when displaying items on the frontend
	apiv1.POST("/list/items", v1.AddItem)
	apiv1.PUT("/items", v1.UpdateItems)
	apiv1.PUT("/item", v1.UpdateItem)
	apiv1.DELETE("/item", v1.DeleteItem)
	apiv1.DELETE("/list/:id", v1.DeleteShoppinglist)
	apiv1.POST("/participant", v1.AddParticipant)
	apiv1.GET("/participant/requests/:email", v1.GetPendingRequests)
	apiv1.POST("/participant/requests", v1.AcceptRequest)
	apiv1.DELETE("/participant/requests", v1.DeleteRequest)

	apiv1.POST("/resetpassword/verifyid", api.VerifyVerificationId)
	apiv1.POST("/resetpassword", api.SendResetPassword)
	apiv1.POST("/resetpassword/changepassword", api.ChangePassword)

	apiv1.GET("/notifications/n/hasunread", notifications_v1.HasUnreadNotifications)
	apiv1.GET("/notifications", notifications_v1.GetNotifications)
	apiv1.GET("/notification/:notification_id", notifications_v1.GetNotification)
	apiv1.POST("/notifications/n/markall", notifications_v1.MarkAllNotificationsAsRead)
	apiv1.DELETE("/notifications", notifications_v1.DeleteNotification)

	apiv1.GET("/backupcodes", api.GetBackupCodes)
	apiv1.POST("/backupcodes", api.VerifyBackupCode)
	apiv1.POST("/backupcodes/regenerate", api.RegenerateCodes)
	apiv1.POST("/backupcodes/generate", api.GenerateCodes)

	apiv1.POST("/twofactorauthentication", api.UpdateTwoFactorAuthentication)
	r.POST("/twofactorauthentication/verify", api.VerifyTwoFactorAuthentication)

	return r
}
