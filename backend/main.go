package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urento/shoppinglist/middleware/ratelimiter"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/cache"
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
	ratelimiter.Setup()
	cache.Setup()
}

//TODO: Create /user route to retrieve all important user information
//TODO: Create /logout route to invalidate old token
//TODO: Fix Test Equal statement (switch expected and actual)
//TODO: Check JWT stuff
//TODO: Implement Transactions to SQL Queries
//TODO: Implement 2FA and auto filout form when SMS comes in
//TODO: Change JWT Library (https://github.com/square/go-jose)

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
