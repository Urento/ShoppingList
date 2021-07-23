package ratelimiter

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
	"github.com/urento/shoppinglist/pkg/e"
)

var rdb *redis.Client
var rrl *redis_rate.Limiter
var ctx = context.Background()

func Setup() {
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

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	rrl = redis_rate.NewLimiter(rdb)
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := rrl.Allow(ctx, "shoppinglist:123", redis_rate.PerMinute(100))
		code := e.SUCCESS
		if err != nil {
			code = e.ERROR_RATE_LIMITER
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    err,
			})
			c.Abort()
			return
		}

		c.Header("Ratelimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			code = e.ERROR_RATELIMIT_TRY_LATER

			seconds := int(res.RetryAfter / time.Second)
			c.Header("Ratelimit-RetryAfter", strconv.Itoa(seconds))

			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    err,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
