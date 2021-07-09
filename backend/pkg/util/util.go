package util

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

func Setup() {
	var err error
	if IsTesting() {
		if util.PROD {
			err = godotenv.Load()
		} else {
			err = godotenv.Load("../.env")
		}
		if err != nil {
			fmt.Print(err.Error())
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
	} else {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		fmt.Print(err.Error())
	}
	return os.Getenv("TESTING") == "true"
}
