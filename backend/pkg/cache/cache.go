package cache

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

const (
	tokenPrefix = "token:"
	emailPrefix = "email:"
)

var rdb *redis.Client

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

	if redisPassword == "testing" {
		rdb = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: redisPassword,
		})
	}
}

func CacheJWT(email, token string) error {
	var ctx = context.Background()
	//86400 = 24 hours in seconds
	if err := rdb.Set(ctx, tokenPrefix+email, token, 86400).Err(); err != nil {
		return err
	}
	if err := rdb.Set(ctx, emailPrefix+token, email, 86400).Err(); err != nil {
		return err
	}
	rdb.Close()
	return nil
}

func GetJWTByEmail(email string) (interface{}, error) {
	val, err := rdb.Get(context.Background(), tokenPrefix+email).Result()
	if err == redis.Nil {
		return nil, errors.New("jwt token not cached")
	} else if err != nil {
		return nil, err
	}

	rdb.Close()
	return val, err
}

func LoadEnv() error {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING {
		err = godotenv.Load()
	} else {
		err = godotenv.Load()
	}

	return err
}
