package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type User struct {
	EMail                   string `json:"e_mail"`
	Username                string `json:"username"`
	Password                string `json:"password"`
	EmailVerified           bool   `json:"email_verified"`
	Rank                    string `json:"rank"`
	TwoFactorAuthentication bool   `json:"two_factor_authentication"`
	IPAddress               string `json:"ip_address"`
	//TODO: Cache Notifications
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
