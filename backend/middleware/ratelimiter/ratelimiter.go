package ratelimiter

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
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

//int64 = remaining
func GetAndUpdateLimit(c *gin.Context) (int64, error) {
	ctx := context.Background()
	var limit int64
	ip := c.ClientIP()

	val, err := rdb.Get(ctx, "ratelimit:"+ip).Result()
	if err != nil || err == redis.Nil {
		//first time sending a request so create a new entry
		err = rdb.Set(ctx, "ratelimit:"+ip, 1, 180*time.Second).Err()
		if err != nil {
			return 0, err
		}
	} else {
		//is already in the database; just incr the old number
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, errors.New("unable to parse the remaining requests from string to int64")
		}

		if v >= 300 {
			return 0, errors.New("limit reached")
		}
		limit = v + 1
		err = rdb.Set(ctx, "ratelimit:"+ip, limit, redis.KeepTTL).Err()
		if err != nil {
			return 0, err
		}
	}

	return limit, nil
}

func ResetLimit(ip string) error {
	err := rdb.Del(context.Background(), "ratelimit:"+ip).Err()
	return err
}
