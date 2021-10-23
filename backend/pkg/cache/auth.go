package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var (
	redisJwtPrefix            = "jwt:"
	failedLoginAttemptsPrefix = "login_attempts:"
	changePasswordPrefix      = "change_password:"
	tokenPrefix               = "token:"
	emailPrefix               = "email:"
	userPrefix                = "user:"
	totpPrefix                = "totp:"
)

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

	return exists == 1, nil
}

func Check(email, token string) (bool, error) {
	ctx := context.Background()

	//check token
	t, err := GetJWTByEmail(email)
	if err != nil {
		//error is probably just that the jwt token is not cached
		return false, nil
	}

	if t != token {
		return false, nil
	}

	//update ttl
	pipe := rdb.Pipeline()

	ttl, err := rdb.TTL(ctx, tokenPrefix+email).Result()
	if err != nil {
		return false, err
	}

	err = pipe.Expire(ctx, tokenPrefix+email, ttl+2*time.Hour).Err()
	if err != nil {
		return false, err
	}

	err = pipe.Expire(ctx, redisJwtPrefix+email, ttl+2*time.Hour).Err()
	if err != nil {
		return false, err
	}

	err = pipe.Expire(ctx, emailPrefix+token, ttl+2*time.Hour).Err()
	if err != nil {
		return false, err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func IsTokenValid(token string) (bool, error) {
	exists, err := rdb.Exists(context.Background(), emailPrefix+token).Result()
	return exists == 1, err
}

func DeleteTokenByEmail(email, token string) (bool, error) {
	ctx := context.Background()
	pipe := rdb.Pipeline()

	exists, err := EmailExists(email)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("email not cached")
	}

	err = pipe.Del(ctx, tokenPrefix+email).Err()
	if err != nil {
		return false, err
	}

	err = pipe.Del(ctx, emailPrefix+token).Err()
	if err != nil {
		return false, err
	}

	_, err = pipe.Exec(ctx)
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

	return kds.SecretId == secretId, nil
}

func HasSecretId(email string) (string, bool, error) {
	ctx := context.Background()

	val, err := rdb.Get(ctx, redisJwtPrefix+email).Result()
	if err != nil || err == redis.Nil {
		log.Print(err)
		return "", false, nil
	}

	var jwtModel JWTModel
	if err := json.Unmarshal([]byte(val), &jwtModel); err != nil {
		return "", false, err
	}

	return jwtModel.SecretId, true, nil
}

func InvalidateSecretId(email string) error {
	err := rdb.Del(context.Background(), redisJwtPrefix+email).Err()
	return err
}

func GetFailedLoginAttempts(ctx context.Context, email string) (int, error) {
	failedAttempts, err := rdb.Get(ctx, failedLoginAttemptsPrefix+email).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return failedAttempts, nil
}

func HasFailedLoginAttempts(ctx context.Context, email string) (bool, error) {
	exists, err := rdb.Exists(ctx, failedLoginAttemptsPrefix+email).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func UpdateFailedLoginAttempts(ctx context.Context, email string) error {
	has, err := HasFailedLoginAttempts(ctx, email)
	if err != nil {
		return err
	}

	if has {
		err = rdb.Incr(ctx, failedLoginAttemptsPrefix+email).Err()
		if err != nil {
			return err
		}

		//reset ttl
		err = rdb.Expire(ctx, failedLoginAttemptsPrefix, 600*time.Second).Err()
		if err != nil {
			return err
		}
	} else {
		err := rdb.Set(ctx, failedLoginAttemptsPrefix+email, 1, 600*time.Second).Err()
		if err != nil {
			return err
		}
	}
	return err
}

func ClearFailedLoginAttempts(ctx context.Context, email string) error {
	err := rdb.Del(ctx, failedLoginAttemptsPrefix+email).Err()
	return err
}

func ActivateResetPassword(ctx context.Context, email string) error {
	err := rdb.Set(ctx, changePasswordPrefix+email, true, 60*time.Minute).Err()
	return err
}

func CanResetPassword(ctx context.Context, email string) (bool, error) {
	r, err := rdb.Get(ctx, changePasswordPrefix+email).Bool()
	if err != nil && err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return r, nil
}

func RemoveResetPassword(ctx context.Context, email string) error {
	err := rdb.Del(ctx, changePasswordPrefix+email).Err()
	return err
}
