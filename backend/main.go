package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urento/shoppinglist/middleware/ratelimiter"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/setting"
	"github.com/urento/shoppinglist/pkg/util"
	routers "github.com/urento/shoppinglist/router"
)

func init() {
	setting.Setup()
	models.Setup()
	util.Setup()
	ratelimiter.Setup()
	cache.Setup()
}

//TODO: Check JWT stuff
//TODO: Implement Transactions to SQL Queries
//TODO: Implement 2FA and auto filout form when SMS comes in
//TODO: Add Update User Route and allow to only update specific fields
//TODO: BeforeEach and AfterEach in Tests?
//TODO: Implement User Cache

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
