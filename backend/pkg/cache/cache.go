package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func CacheJWT(email, token string) error {
	ctx := context.Background()
	t := 24 * time.Hour

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.Set(ctx, tokenPrefix+email, token, t).Err()
		if err != nil {
			return err
		}

		err = pipe.Set(ctx, emailPrefix+token, email, t).Err()
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

func InvalidateSpecificJWTToken(email, token string) error {
	ctx := context.Background()
	pipe := rdb.Pipeline()

	err := pipe.Del(ctx, emailPrefix+token).Err()
	if err != nil {
		return err
	}

	err = pipe.Del(ctx, redisJwtPrefix+email).Err()
	if err != nil {
		return err
	}

	err = pipe.Del(ctx, tokenPrefix+email).Err()
	if err != nil {
		return err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DoesTokenBelongToEmail(email, token string) (bool, error) {
	val, err := rdb.Get(context.Background(), tokenPrefix+email).Result()
	if err != nil {
		return false, err
	}

	if val != token {
		return false, nil
	}

	return true, nil
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
	ctx := context.Background()
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

type User struct {
	EMail                   string `json:"e_mail"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	EmailVerified           bool   `json:"email_verified"`
	Rank                    string `json:"rank"`
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
	IPAddress               string `json:"ip_address"`
}

func (user User) CacheUser() error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rdb.Set(context.Background(), userPrefix+user.EMail, b, 0).Err()
	if err != nil {
		return err
	}

	return err
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

func GetTwoFactorAuthenticationStatus(email string) (bool, error) {
	user, err := GetUser(email)
	if err != nil {
		return false, err
	}

	return user.TwoFactorAuthentication, nil
}

//TODO: also let only specific fields get updated
func UpdateUser(user User) error {
	ctx := context.Background()
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.Del(ctx, userPrefix+user.EMail).Err()
		if err != nil {
			return err
		}

		b, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = pipe.Set(ctx, userPrefix+user.EMail, b, 0).Err()
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

type JWTModel struct {
	SecretId string `json:"secret_id"`
	Email    string `json:"email"`
}

func GenerateSecretId(email string) (string, error) {
	existingSecretId, has, err := HasSecretId(email)
	if err != nil {
		return "", err
	}

	if has {
		return existingSecretId, nil
	}

	secretId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	jwtModel := JWTModel{
		Email:    email,
		SecretId: secretId.String(),
	}

	b, err := json.Marshal(jwtModel)
	if err != nil {
		return "", err
	}

	err = rdb.Set(context.Background(), redisJwtPrefix+email, b, 86400*time.Second).Err()
	if err != nil {
		return "", err
	}

	return secretId.String(), nil
}

func VerifySecretId(email, secretId string) (bool, error) {
	ctx := context.Background()

	//verify secretId
	obj, err := rdb.Get(ctx, redisJwtPrefix+email).Result()
	if err != nil || err == redis.Nil {
		return false, errors.New("secretid is not valid")
	}

	var kds JWTModel
	if err := json.Unmarshal([]byte(obj), &kds); err != nil {
		return false, err
	}

	if kds.SecretId == secretId {
		return true, nil
	}

	return false, nil
}

func HasSecretId(email string) (string, bool, error) {
	ctx := context.Background()

	val, err := rdb.Get(ctx, redisJwtPrefix+email).Result()
	if err != nil || err == redis.Nil {
		log.Print(err)
		return "", false, nil
	}

	var kds JWTModel
	if err := json.Unmarshal([]byte(val), &kds); err != nil {
		return "", false, err
	}

	return kds.SecretId, true, nil
}

func InvalidateSecretId(email string) error {
	err := rdb.Del(context.Background(), redisJwtPrefix+email).Err()
	return err
}

func CacheTOTPSecret(email, secret string) error {
	err := rdb.Set(context.Background(), totpPrefix+email, secret, 0).Err()
	return err
}

func GetTOTPSecret(email string) (string, error) {
	val, err := rdb.Get(context.Background(), totpPrefix+email).Result()
	if err != nil {
		return "not cached", errors.New("totp secret is not cached")
	}
	return val, nil
}

func DeleteTOTPSecret(email string) error {
	err := rdb.Del(context.Background(), totpPrefix+email).Err()
	return err
}

func IsTOTPSecretCached(email string) (bool, error) {
	exists, err := rdb.Exists(context.Background(), totpPrefix+email).Result()
	return exists == 1, err
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
