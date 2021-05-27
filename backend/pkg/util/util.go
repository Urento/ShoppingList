package util

import "github.com/urento/shoppinglist/pkg/setting"

func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
