package totp

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/xlzd/gotp"
)

func Disable(email string, appGin *app.Gin) {
	err := models.SetTwoFactorAuthentication(email, false)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_CHANING_TWOFACTORAUTHENTICATION_STATUS, map[string]string{"success": "false"})
		return
	}

	err = cache.DeleteTOTPSecret(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_CHANING_TWOFACTORAUTHENTICATION_STATUS, map[string]string{"success": "false"})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true", "verified": "true"})
}

func Enable(email string, appGin *app.Gin) []byte {
	key, err := totp.Generate(totp.GenerateOpts{
		AccountName: email,
		Issuer:      "Shoppinglist",
	})
	if err != nil {
		return []byte(err.Error())
	}

	err = cache.CacheTOTPSecret(email, key.Secret())
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_CHANING_TWOFACTORAUTHENTICATION_STATUS, nil)
		return []byte(err.Error())
	}

	//Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_CONVERTING_KEY_TO_IMG, nil)
		return []byte(err.Error())
	}

	png.Encode(&buf, img)
	return buf.Bytes()
}

func GetQRCodeBase64String(email string, data []byte) string {
	imgBase64Str := base64.StdEncoding.EncodeToString(data)
	//img2html := "<img src=\"data:image/png;base64," + imgBase64Str + "\" />"
	return imgBase64Str
}

func Verify(email, a string, enableAfter bool) (bool, error) {
	timestamp := time.Now().Unix()

	secret, err := cache.GetTOTPSecret(email)
	if err != nil {
		return false, err
	}

	if enableAfter {
		err := models.SetTwoFactorAuthentication(email, true)
		if err != nil {
			log.Print(err)
			return false, err
		}
	}

	otp := gotp.NewDefaultTOTP(secret)
	return otp.Verify(a, int(timestamp)), nil
}
