package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alexedwards/argon2id"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/totp"
	"github.com/urento/shoppinglist/pkg/util"
)

type Auth struct {
	Email                   string `valid:"Required;"`
	Username                string
	Password                string `valid:"Required"`
	EmailVerified           bool
	TwoFactorAuthentication bool
}

func Check(c *gin.Context) {
	appGin := app.Gin{C: c}

	token, err := GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{"success": "false"})
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{"success": "false"})
		return
	}

	check, err := cache.Check(email, token)
	if !check || err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_TOKEN_INVALID, map[string]string{"success": "false"})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}

func GetUser(c *gin.Context) {
	appGin := app.Gin{C: c}

	token, err := GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
		return
	}

	data, err := models.GetUser(email)
	if err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_RETRIEVING_USER_DATA, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, data)
}

type UpdateUserStruct struct {
	EMail         string `json:"e_mail"`
	EmailVerified bool   `json:"email_verified"`
	Username      string `json:"username"`
	OldPassword   string `json:"old_password"`
	Password      string `json:"password"`
	WithPassword  bool   `json:"with_password"`
}

func UpdateUser(c *gin.Context) {
	appGin := app.Gin{C: c}

	token, err := GetCookie(c)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	var data UpdateUserStruct

	if err := c.BindJSON(&data); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	b, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_BINDING_JSON_DATA, map[string]string{
			"success": "false",
			"message": "error decoding struct",
		})
		return
	}

	var lokifdgh models.Auth
	if err := json.Unmarshal([]byte(b), &lokifdgh); err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_BINDING_JSON_DATA, map[string]string{
			"success": "false",
			"message": "error decoding struct 2",
		})
		return
	}

	if data.WithPassword {
		//TODO: maybe even check the cache and not postgres
		ok, err := models.CheckPassword(email, data.OldPassword)
		if !ok || err != nil {
			appGin.Response(http.StatusBadRequest, e.ERROR_WRONG_OLD_PASSWORD, map[string]string{
				"success": "false",
				"message": "Old Password is incorrect!",
			})
			return
		}

		passwordHash, err := argon2id.CreateHash(lokifdgh.Password, argon2id.DefaultParams)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.ERROR_ENCRYPTING_PASSWORD, map[string]string{
				"success": "false",
				"message": "error encrypting password",
			})
			return
		}

		lokifdgh.Password = passwordHash
	}

	err = lokifdgh.UpdateUser(email)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_UPDATING_USER, map[string]string{
			"success": "false",
			"message": "error updating user",
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	//TODO: Maybe decode the token and check expire time

	var user LoginUser

	if err := c.BindJSON(&user); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email := user.Email
	password := user.Password
	ip := c.ClientIP()

	a := Auth{Email: email, Password: password}
	ok, err := valid.Valid(&a)
	if !ok {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//check if the user still has a valid token
	/*ex, err := cache.EmailExists(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}*/

	//TODO: Maybe decode token and validate the actual token again

	//user still has a valid token
	/*if ex {
		token, err := cache.GetJWTByEmail(email)
		if err != nil {
			log.Print(err)
			appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			return
		}

		err = SetCookie(c, token)
		if err != nil {
			log.Print(err)
			appGin.Response(http.StatusInternalServerError, e.ERROR_SETTING_SESSION_TOKEN, map[string]string{
				"error":   err.Error(),
				"success": "false",
			})
			return
		}

		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"token": token,
		})
		return
	}*/

	// if the user doesnt have a valid token in cache = generate new one
	exists, err := models.CheckAuth(email, password, ip)
	if err != nil && err.Error() == "too many failed login attempts" {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]string{
			"success": "false",
			"error":   "too many failed login attempts",
		})
		return
	}
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]string{
			"success": "false",
			"error":   "wrong email or password",
		})
		return
	}

	if !exists {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH, map[string]string{
			"success": "false",
			"error":   "wrong email or password",
		})
		return
	}

	has, err := cache.IsTOTPSecretCached(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusUnauthorized, e.ERROR_CHECKING_IF_TOTP_IS_ENABLED, map[string]string{
			"success": "false",
			"error":   "error while getting totp",
		})
		return
	}

	if has {
		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success": "true",
			"otp":     "true",
		})
		return
	}

	token, err := util.GenerateToken(email, false)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, map[string]string{
			"success": "false",
			"error":   "error while generating jwt",
		})
		return
	}

	err = SetCookie(c, token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_SETTING_SESSION_TOKEN, map[string]string{
			"error":   err.Error(),
			"success": "false",
			"otp":     "false",
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token":   token,
		"success": "true",
		"totp":    "false",
	})
}

type LogoutSettings struct {
	LogoutEveryone bool `json:"logout_everyone"`
}

func Logout(c *gin.Context) {
	appGin := app.Gin{C: c}

	token, err := GetCookie(c)
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"success": "false",
			"message": "You have to be logged in to log out!",
		})
		return
	}

	var logoutSettings LogoutSettings

	if err := c.BindJSON(&logoutSettings); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"message": "jwt token is invalid",
		})
		return
	}

	//if you want to invalidate all jwt tokens and log everyone out
	if logoutSettings.LogoutEveryone {
		/*err := cache.InvalidateAllJWTTokens(email)
		if err != nil {
			log.Print(err)
			appGin.Response(http.StatusInternalServerError, e.ERROR_INVALIDATING_JWT_TOKENS, map[string]string{
				"success": "false",
				"message": "error while invalidating jwt token",
			})
		}*/
		//TODO

		RemoveCookie(c)

		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success": "true",
		})
		return
	}

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"success": "false",
		})
		return
	}

	//if its just a normal logout
	ok, err := cache.DeleteTokenByEmail(email, token)
	if err != nil || !ok {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_WHILE_INVALIDATING_TOKEN, map[string]string{
			"success": strconv.FormatBool(ok),
			"error":   "error while invalidating token",
		})
	}

	RemoveCookie(c)

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
	})
}

type InvalidateSpecificJWTTokenStruct struct {
	JWTToken string `json:"jwt_token"`
}

func InvalidateSpecificJWTToken(c *gin.Context) {
	appGin := app.Gin{C: c}

	token, err := GetCookie(c)
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"success": "false",
			"message": "You have to be logged in to log out!",
		})
		return
	}

	var jwtTokenSettings InvalidateSpecificJWTTokenStruct

	if err := c.BindJSON(&jwtTokenSettings); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"message": "jwt token is invalid",
		})
		return
	}

	ok, err := cache.DoesTokenBelongToEmail(email, jwtTokenSettings.JWTToken)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_JWT_TOKEN_DOES_NOT_BELONG_TO_EMAIL, map[string]string{
			"success": "false",
			"message": "JWT Token does not belong to your account!",
		})
		return
	}

	if !ok {
		appGin.Response(http.StatusBadRequest, e.ERROR_JWT_TOKEN_DOES_NOT_BELONG_TO_EMAIL, map[string]string{
			"success": "false",
			"message": "JWT Token does not belong to your account!",
		})
		return
	}

	err = cache.InvalidateSpecificJWTToken(email, jwtTokenSettings.JWTToken)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_INVALIDATING_JWT_TOKENS, map[string]string{
			"success": "false",
			"message": "An error occurred while invalidating a specific JWT Token!",
		})
		return
	}

	appGin.Response(http.StatusBadRequest, e.ERROR_JWT_TOKEN_DOES_NOT_BELONG_TO_EMAIL, map[string]string{
		"success": "true",
	})
}

func RemoveCookie(ctx *gin.Context) {
	domain := os.Getenv("DOMAIN")
	ctx.SetCookie("token", "", -1, "/", domain, util.IsProd(), true)
}

type RegisterUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateAccount(c *gin.Context) {
	appGin := app.Gin{C: c}
	valid := validation.Validation{}

	var user RegisterUser

	if err := c.BindJSON(&user); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email := user.Email
	username := user.Username
	password := user.Password

	if len(username) > 32 {
		appGin.Response(http.StatusBadRequest, e.ERROR_USERNAME_TOO_LONG, map[string]string{
			"success": "false",
			"error":   "username can not be longer than 32 characters",
		})
		return
	}

	ip := c.ClientIP()

	a := Auth{Email: email, Username: username, Password: password}
	ok, _ := valid.Valid(a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "validation error",
		})
		return
	}

	err := models.CreateAccount(email, username, password, ip)
	if err != nil {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.ERROR_CREATING_ACCOUNT, map[string]string{
			"error":   "email is already being used",
			"success": "false",
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"created":  "true",
		"email":    email,
		"username": username,
	})
}

type EmailVerification struct {
	SecretID string `json:"secret_id"`
	Email    string `json:"email"`
}

func UpdateEmailVerified(c *gin.Context) {
	//TODO
}

type TwoFactorAuthentictionUpdate struct {
	OTP    string `json:"otp"`
	Status bool   `json:"status"`
}

func UpdateTwoFactorAuthentication(c *gin.Context) {
	appGin := app.Gin{C: c}
	token, err := GetCookie(c)
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"success":  "false",
			"message":  "You have to be logged in to log out!",
			"verified": "false",
		})
		return
	}

	var data TwoFactorAuthentictionUpdate

	if err := c.BindJSON(&data); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, map[string]string{"success": "false", "verified": "false"})
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{"success": "false", "verified": "false"})
		return
	}

	//check if the two factor authentication status is the same as in the cache; if yes = dont process the request
	//TODO: UNCOMMENT ONCE I IMPLEMENT USER CACHING
	/*currentStatus, err := cache.GetTwoFactorAuthenticationStatus(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_TWOFACTORAUTHENTICATION_STATUS_FROM_CACHE, nil)
		return
	}

	if currentStatus == data.Status {
		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success": "false",
			"message": "can't update the status if the status in the request is the same as in the cache",
		})
		return
	}*/

	enabled, err := models.IsTwoFactorEnabled(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_TWOFACTORAUTHENTICATION_STATUS_FROM_CACHE, map[string]string{"success": "false", "verified": "false"})
		return
	}

	if enabled == data.Status {
		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success":  "false",
			"message":  "You can't change your status right now! Try again later!",
			"verified": "false",
		})
		return
	}

	if !enabled {
		bytes := totp.Enable(email, &appGin)
		qrCode := totp.GetQRCodeBase64String(email, bytes)

		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success":  "true",
			"verified": "true",
			"img":      qrCode,
		})
		return
	}

	ok, err := totp.Verify(email, data.OTP, false)
	if err != nil || !ok {
		appGin.Response(http.StatusBadRequest, e.ERROR_VERIFYING_OTP, map[string]string{"success": "false", "message": "OTP was wrong"})
		return
	}

	totp.Disable(email, &appGin)
}

type VerifyTOTP struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	LoginAfter  bool   `json:"login_after"`
	EnableAfter bool   `json:"enable_after"`
}

func VerifyTwoFactorAuthentication(c *gin.Context) {
	appGin := app.Gin{C: c}

	var data VerifyTOTP

	if err := c.BindJSON(&data); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, map[string]string{"success": "false", "verified": "false"})
		return
	}

	email := data.Email

	enabled, err := models.IsTwoFactorEnabled(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_TWOFACTORAUTHENTICATION_STATUS_FROM_CACHE, map[string]string{"success": "false", "verified": "false"})
		return
	}

	if !enabled && !data.EnableAfter {
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_TWOFACTORAUTHENTICATION_STATUS_FROM_CACHE, map[string]string{
			"success":  "false",
			"message":  "TOTP not activated!",
			"verified": "false",
		})
		return
	}

	ok, err := totp.Verify(email, data.OTP, data.EnableAfter)
	if err != nil || !ok {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_VERIFYING_OTP, map[string]string{"success": "false", "verified": "false"})
		return
	}

	if data.LoginAfter && ok {
		token, err := util.GenerateToken(email, false)
		if err != nil {
			appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
			return
		}

		err = SetCookie(c, token)
		if err != nil {
			log.Print(err)
			appGin.Response(http.StatusInternalServerError, e.ERROR_SETTING_SESSION_TOKEN, map[string]string{
				"error":    err.Error(),
				"token":    token,
				"success":  "false",
				"verified": "false",
			})
			return
		}
		appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
			"success":  "true",
			"verified": "true",
			"token":    token,
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true", "verified": "true"})
}

type ResetPasswordRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}

func ResetPasswordFromUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var form ResetPasswordRequest

	if err := c.BindJSON(&form); err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, map[string]string{"success": "false", "verified": "false"})
		return
	}

	token, err := util.GetCookie(c)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusBadRequest, e.ERROR_GETTING_HTTPONLY_COOKIE, map[string]string{
			"error":   err.Error(),
			"success": "false",
		})
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
		})
		return
	}

	err = models.ResetPasswordFromUser(email, form.Password, form.OldPassword, true)
	if err != nil {
		log.Print(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
			"error":   "error while resetting password from user",
		})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{"success": "true"})
}

func SetCookie(ctx *gin.Context, token string) error {
	domain := os.Getenv("DOMAIN")
	ctx.SetCookie("token", token, 24*60*60, "/", domain, util.IsProd(), true)
	return nil
}

func GetCookie(ctx *gin.Context) (string, error) {
	token, err := ctx.Request.Cookie("token")
	if err != nil {
		return "", err
	}

	if len(token.Value) <= 0 {
		return "", errors.New("cookie 'token' has to be longer than 0 charcters")
	}

	if len(token.Value) <= 50 {
		return "", errors.New("cookie 'token' has to be longer than 50 charcters")
	}

	return token.Value, nil
}
