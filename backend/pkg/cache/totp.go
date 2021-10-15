package cache

import (
	"context"
	"errors"
)

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
