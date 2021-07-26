package cache

import (
	"context"
	"encoding/json"
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
	userPrefix  = "user:"
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

type User struct {
	EMail                   string `json:"e_mail"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	EmailVerified           bool   `json:"email_verified"`
	Rank                    string `json:"rank"`
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
}

func CacheUser(user User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rdb.Set(context.Background(), userPrefix+user.EMail, b, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUser(email string) (*User, error) {
	var user User
	u, err := rdb.Get(context.Background(), userPrefix+email).Result()
	if err != nil {
		return nil, err
	}

	d := []byte(u)
	err = json.Unmarshal(d, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//TODO: Let only specific fields get updated
func UpdateUser(user User) error {
	var ctx = context.Background()
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := rdb.Del(ctx, userPrefix+user.EMail).Err()
		if err != nil {
			return err
		}

		b, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = rdb.Set(ctx, userPrefix+user.EMail, b, 0).Err()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(email string) error {
	err := rdb.Del(context.Background(), userPrefix+email).Err()
	return err
}

func IsUserCached(email string) (bool, error) {
	exists, err := rdb.Exists(context.Background(), userPrefix+email).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func CacheJWT(email, token string) error {
	var ctx = context.Background()
	//86400 = 24 hours in seconds

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := rdb.Set(ctx, tokenPrefix+email, token, 86400*time.Second).Err()
		if err != nil {
			return err
		}

		err = rdb.Set(ctx, emailPrefix+token, email, 86400*time.Second).Err()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func GetJWTByEmail(email string) (string, error) {
	val, err := rdb.Get(context.Background(), tokenPrefix+email).Result()
	if err == redis.Nil {
		return "", errors.New("jwt token not cached")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func GetEmailByJWT(token string) (string, error) {
	val, err := rdb.Get(context.Background(), emailPrefix+token).Result()
	if err == redis.Nil {
		return "", errors.New("jwt token not cached")
	} else if err != nil {
		return "", err
	}
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

func Check(email, token string) (bool, error) {
	exists, err := EmailExists(email)
	if err != nil || !exists {
		return false, err
	}

	t, err := GetJWTByEmail(email)
	if err != nil {
		return false, err
	}

	if t != token {
		return false, nil
	}

	return true, nil
}

func IsTokenValid(token string) (bool, error) {
	exists, err := rdb.Exists(context.Background(), emailPrefix+token).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func DeleteTokenByEmail(email, token string) (bool, error) {
	var ctx = context.Background()
	exists, err := EmailExists(email)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("email not cached")
	}

	err = rdb.Del(ctx, tokenPrefix+email).Err()
	if err != nil {
		return false, err
	}

	err = rdb.Del(ctx, emailPrefix+token).Err()
	if err != nil {
		return false, err
	}

	return true, nil
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
