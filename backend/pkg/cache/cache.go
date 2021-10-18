package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

var (
	redisJwtPrefix = "jwt:"
	tokenPrefix    = "token:"
	emailPrefix    = "email:"
	userPrefix     = "user:"
	totpPrefix     = "totp:"
)

var rdb *redis.Client

//TODO: Add Cache for Shoppinglists

func Setup(auth_test bool) {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if auth_test {
		err = godotenv.Load("../.env")
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING && !auth_test {
		err = godotenv.Load("../../.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		panic(err)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddr := os.Getenv("REDIS_ADDR")

	if redisPassword == "testing" {
		rdb = redis.NewClient(&redis.Options{
			Addr: redisAddr,
			DB:   0,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})
	}
}
