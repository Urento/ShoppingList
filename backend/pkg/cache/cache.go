package cache

import (
	"log"
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

func LoadEnv() error {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING {
		err = godotenv.Load("../.env")
	} else {
		err = godotenv.Load()
	}
	return err
}
