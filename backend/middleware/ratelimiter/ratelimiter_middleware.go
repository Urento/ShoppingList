package ratelimiter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Ratelimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		remaining, err := GetAndUpdateLimit(c)
		log.Print(remaining)

		//TODO: make exceptions for some routes maybe?

		c.Header("Ratelimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("Ratelimit-Limit", "300")

		if err != nil && err.Error() == "limit reached" {
			log.Print(err)
			c.JSON(http.StatusTooManyRequests, Response{
				Error:   "Ratelimit reached!",
				Message: "Try again later!",
			})
			c.Abort()
			return
		}

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, Response{
				Error:   err.Error(),
				Message: "We are currently unable to process your requests! Try again later!",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
