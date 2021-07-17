package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/logging"
	"github.com/urento/shoppinglist/pkg/setting"
	"github.com/urento/shoppinglist/pkg/util"
	routers "github.com/urento/shoppinglist/router"
)

func init() {
	setting.Setup()
	logging.Setup()
	models.Setup()
	util.Setup()
}

//TODO: Create /user route to retrieve all important user information

func main() {
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("Listening on Port %s", endPoint)

	server.ListenAndServe()
}
