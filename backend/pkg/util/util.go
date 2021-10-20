package util

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

const letterAndNumberBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func Setup() {
	var err error
	if IsTesting() {
		if util.PROD {
			err = godotenv.Load()
		} else if util.GITHUB_TESTING {
			err = nil
		} else if util.LOCAL_TESTING {
			err = godotenv.Load("../.env")
		} else {
			err = godotenv.Load()
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	jwtSecret = []byte(os.Getenv("JwtSecret"))
}

func IsTesting() bool {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING {
		err = godotenv.Load("../.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		fmt.Print(err.Error())
	}

	return os.Getenv("TESTING") == "true"
}

func IsProd() bool {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else if util.GITHUB_TESTING {
		err = nil
	} else if util.LOCAL_TESTING {
		err = godotenv.Load("../.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		fmt.Print(err.Error())
	}

	return os.Getenv("ENVIRONMENT") == "production"
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

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterAndNumberBytes[rand.Intn(len(letterAndNumberBytes))]
	}
	return string(b)
}

func RandomEmail() string {
	return StringWithCharset(32) + "@gmail.com"
}

func StringWithCharset(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letterAndNumberBytes[seededRand.Intn(len(letterAndNumberBytes))]
	}
	return string(b)
}

func RandomInt() int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int()
	return n
}

func RandomIntWithLength(length int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(length)
	return n
}

func RandomPosition() int64 {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(seededRand.Intn(math.MaxInt64))
}

func RandomBoolean() bool {
	s := rand.Intn(2)
	return s > 1
}

func RandomRank() string {
	nmb := rand.Intn(2)
	if nmb <= 1 {
		return "admin"
	}
	return "default"
}

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func StringArrayToArray(before []string, i int) string {
	replacer := strings.NewReplacer("{", "", "}", "")
	output := replacer.Replace(before[0])

	s := strings.Split(output, ",")
	return s[i]
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
