package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs details about each incoming HTTP request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Process request
		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		log.Printf("%s %s - %d (%s)", method, path, statusCode, latency)
	}
}
