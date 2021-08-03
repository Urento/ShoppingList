package ratelimiter

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRatelimitMiddleware(t *testing.T) {
	Setup()

	r := gin.Default()

	r.Use(Ratelimiter())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	go r.Run(":9999")

	c := &http.Client{}
	for i := 0; i < 302; i++ {
		resp, err := c.Get("http://127.0.0.1:9999")
		if err != nil {
			t.Errorf("An error occurred while making the requests: %s", err)
		}

		limitHeader := resp.Header.Get("Ratelimit-Limit")

		if limitHeader != "300" {
			t.Error("Ratelimit-Limit header not set!")
		}

		if i == 301 && resp.StatusCode != 429 {
			t.Error("Ratelimit not detected")
		}
	}

	err := ResetLimit("127.0.0.1")
	if err != nil {
		t.Errorf("Error while resetting limit: %s", err)
	}
}
