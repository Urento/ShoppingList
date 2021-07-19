package cache

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

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
			DB:   0,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: redisPassword,
			DB:       0,
		})
	}
}

func CacheJWT(email, token string) error {
	var ctx = context.Background()
	var err error
	//86400 = 24 hours in seconds
	err = rdb.Set(ctx, tokenPrefix+email, token, 0).Err()
	rdb.Do(context.Background(), "EXPIRE", tokenPrefix+email, 86400)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, emailPrefix+token, email, 0).Err()
	rdb.Do(context.Background(), "EXPIRE", emailPrefix+token, 86400)
	if err != nil {
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
	return val, nil
}

func GetEmailByJWT(token string) (string, error) {
	val, err := rdb.Get(context.Background(), emailPrefix+token).Result()
	if err == redis.Nil {
		return "", errors.New("jwt token not cached")
	} else if err != nil {
		return "", err
	}

	rdb.Close()
	return val, nil
}

func EmailExists(email string) (bool, error) {
	exists, err := rdb.Exists(context.Background(), tokenPrefix+email).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func GetTTLByEmail(email string) (time.Duration, error) {
	ttl, err := rdb.TTL(context.Background(), tokenPrefix+email).Result()
	if err != nil {
		return -1, err
	}
	return ttl, nil
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
