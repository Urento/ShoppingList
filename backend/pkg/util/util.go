package util

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

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
