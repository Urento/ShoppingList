package api

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/alexedwards/argon2id"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/app"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/e"
	"github.com/urento/shoppinglist/pkg/util"
	auth "github.com/urento/shoppinglist/services"
)

var (
	sessionName = "token"
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
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"success": "false",
		})
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil || len(email) <= 0 {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_EMAIL_BY_JWT, map[string]string{
			"success": "false",
		})
		return
	}

	check, err := cache.Check(email, token)
	if !check || err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_TOKEN_INVALID, map[string]string{
			"success": "false",
		})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"token":   token,
	})
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

	authService := auth.Auth{EMail: email}
	data, err := authService.GetUser()
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
	ip, err := GetClientIPHelper(c.Request)
	if err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_IP, map[string]string{
			"success": "false",
			"error":   "error while resolving ip address",
		})
	}

	a := Auth{Email: email, Password: password}
	ok, err := valid.Valid(&a)
	if !ok {
		log.Print(err)
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	//check if the user still has a valid token
	ex, err := cache.EmailExists(email)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	//TODO: Maybe decode token and validate the actual token again

	//user still has a valid token
	if ex {
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
	}

	// if the user doesnt have a valid token in cache = generate new one
	authService := auth.Auth{EMail: email, Password: password, IPAddress: ip}
	exists, err := authService.Check()
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !exists {
		appGin.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(email)
	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
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
		"token":   token,
		"success": "true",
	})
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

	if len(token) <= 0 {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"success": "false",
		})
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

	ip, err := GetClientIPHelper(c.Request)
	if err != nil {
		appGin.Response(http.StatusBadRequest, e.ERROR_GETTING_IP, map[string]string{
			"success": "false",
			"error":   "error resolving the ip address",
		})
		return
	}

	a := Auth{Email: email, Username: username, Password: password}
	ok, _ := valid.Valid(a)
	if !ok {
		app.MarkErrors(valid.Errors)
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]string{
			"error": "validation error",
		})
		return
	}

	acc := auth.Auth{EMail: email, Username: username, Password: password, EmailVerified: false, Rank: "default", IPAddress: ip}
	err = acc.Create()
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

func UpdateEmailVerified(c *gin.Context) {
	//TODO
}

func DisplayTwoFactorAuthenticationQRCode(c *gin.Context) {

}

type TwoFactorAuthentictionUpdate struct {
	Status    bool   `json:"status"`
	UserToken string `json:"user_token"`
}

func UpdateTwoFactorAuthentication(c *gin.Context) {
	//TODO https://github.com/sec51/twofactor
	appGin := app.Gin{C: c}
	token, err := GetCookie(c)
	if err != nil {
		appGin.Response(http.StatusUnauthorized, e.ERROR_NOT_AUTHORIZED, map[string]string{
			"success": "false",
			"message": "You have to be logged in to log out!",
		})
		return
	}

	var data TwoFactorAuthentictionUpdate

	if err := c.BindJSON(&data); err != nil {
		log.Print(err)
		appGin.Response(http.StatusBadRequest, e.ERROR_BINDING_JSON_DATA, nil)
		return
	}

	email, err := cache.GetEmailByJWT(token)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_GETTING_EMAIL_BY_JWT, nil)
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

	err = models.SetTwoFactorAuthentication(email, data.Status)
	if err != nil {
		log.Print(err)
		appGin.Response(http.StatusInternalServerError, e.ERROR_CHANING_TWOFACTORAUTHENTICATION_STATUS, nil)
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"success": "true",
		"status":  strconv.FormatBool(data.Status),
	})
}

// GetClientIPHelper gets the client IP using a mixture of techniques.
// This is how it is with golang at the moment.
func GetClientIPHelper(req *http.Request) (ipResult string, errResult error) {

	// Try lots of ways :) Order is important.

	//  Try Request Header ("Origin")
	url, err := url.Parse(req.Header.Get("Origin"))
	if err == nil {
		host := url.Host
		ip, _, err := net.SplitHostPort(host)
		if err == nil {
			log.Printf("debug: Found IP using Header (Origin) sniffing. ip: %v", ip)
			return ip, nil
		}
	}

	// Try by Request
	ip, err := getClientIPByRequestRemoteAddr(req)
	if err == nil {
		log.Printf("debug: Found IP using Request sniffing. ip: %v", ip)
		return ip, nil
	}

	// Try Request Headers (X-Forwarder). Client could be behind a Proxy
	ip, err = getClientIPByHeaders(req)
	if err == nil {
		log.Printf("debug: Found IP using Request Headers sniffing. ip: %v", ip)
		return ip, nil
	}

	err = errors.New("error: Could not find clients IP address")
	return "", err
}

// getClientIPByRequest tries to get directly from the Request.
// https://blog.golang.org/context/userip/userip.go
func getClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {

	// Try via request
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("debug: Getting req.RemoteAddr %v", err)
		return "", err
	} else {
		log.Printf("debug: With req.RemoteAddr found IP:%v; Port: %v", ip, port)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		message := "parsing ip from request.remoteaddr got nothing"
		log.Print(message)
		return "", errors.New(message)

	}
	log.Printf("debug: Found IP: %v", userIP)
	return userIP.String(), nil

}

// getClientIPByHeaders tries to get directly from the Request Headers.
// This is only way when the client is behind a Proxy.
func getClientIPByHeaders(req *http.Request) (ip string, err error) {

	// Client could be behid a Proxy, so Try Request Headers (X-Forwarder)
	ipSlice := []string{}

	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		log.Printf("debug: client request header check gives ip: %v", v)
		if v != "" {
			return v, nil
		}
	}
	err = errors.New("error: Could not find clients IP address from the Request Headers")
	return "", err

}
