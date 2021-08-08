package app

import (
	"log"

	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Print(err.Key, err.Message)
	}
}
