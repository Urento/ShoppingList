package routers

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/urento/shoppinglist/middleware/jwt"
	"github.com/urento/shoppinglist/middleware/ratelimiter"
	util "github.com/urento/shoppinglist/pkg"
	"github.com/urento/shoppinglist/router/api/v1"
	v1 "github.com/urento/shoppinglist/router/api/v1/shoppinglist"
)

type RedisData struct {
	Address  string
	Password string
}

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(ratelimiter.RateLimit())

	redisAddress, redisPassword, err := GetRedisData()
	if err != nil {
		log.Fatal(err)
	}
	//sessionSecret := GetSessionSecret()
	sessionSecret := GetSessionSecret()
	store, _ := redis.NewStore(10, "tcp", redisAddress, redisPassword, []byte(sessionSecret))
	r.Use(sessions.Sessions("loginsession", store))

	//cors setup
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization") //TODO: Maybe remove since I handle logins different now
	r.Use(cors.New(config))

	r.POST("/api/auth", api.Login)
	r.POST("/api/auth/register", api.CreateAccount)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())

	apiv1.POST("/auth/check", api.Check)
	apiv1.GET("/auth/user", api.GetUser)
	apiv1.POST("/auth/logout", api.Logout)

	apiv1.GET("/lists", v1.GetShoppinglists)
	apiv1.POST("/list", v1.CreateShoppinglist)
	apiv1.PUT("/list", v1.EditShoppinglist)
	apiv1.GET("/list/:id", v1.GetShoppinglist)
	apiv1.DELETE("/list/:id", v1.DeleteShoppinglist)

	apiv1.POST("/resetpassword/verifyid", api.VerifyVerificationId)
	apiv1.POST("/resetpassword", api.SendResetPassword)
	apiv1.POST("/resetpassword/changepassword", api.ChangePassword)

	return r
}

func GetRedisData() (string, string, error) {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING {
		err = godotenv.Load("../../.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal(err)
	}

	address := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")

	return address, password, nil
}

func GetSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	return secret
}
