package main

import (
	"log"
	"net/http"

	"github.com/urento/shoppinglist/middleware/ratelimiter"
	"github.com/urento/shoppinglist/models"
	"github.com/urento/shoppinglist/pkg/cache"
	"github.com/urento/shoppinglist/pkg/util"
	routers "github.com/urento/shoppinglist/router"
)

func init() {
	models.Setup()
	util.Setup()
	ratelimiter.Setup()
	cache.Setup(false)
}

//TODO: Check JWT stuff
//TODO: Implement 2FA and auto filout form when SMS comes in
//TODO: BeforeEach and AfterEach in Tests?
//TODO: Implement User Cache
//TODO: Revalidate JWT Token when invalid

func main() {
	routersInit := routers.InitRouter()
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           ":8080",
		Handler:        routersInit,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("Listening on port: %s", ":8080")

	server.ListenAndServe()
}
