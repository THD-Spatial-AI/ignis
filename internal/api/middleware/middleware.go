package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const maxRequestBodyBytes = 1024 // 1 KB — sufficient for {"A_ref": 150.0}

// RequestBodyLimit rejects bodies larger than maxRequestBodyBytes.
// Protects the calculate endpoint from oversized payloads.
func RequestBodyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxRequestBodyBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request body too large",
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodyBytes)
		c.Next()
	}
}

// RequestLogger logs details about each incoming HTTP request.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		log.Printf("%s %s - %d (%s)", method, path, statusCode, latency)
	}
}

// CORS sets CORS headers for browser clients.
//
// Allowed origins are read from the ALLOWED_ORIGINS environment variable as a
// comma-separated list (e.g. "http://localhost:5173,https://app.example.com").
// If the variable is unset or empty, all cross-origin requests are rejected.
// Set ALLOWED_ORIGINS in your .env file — see .env.example.
func CORS() gin.HandlerFunc {
	var allowed []string
	for _, o := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
		if t := strings.TrimSpace(o); t != "" {
			allowed = append(allowed, t)
		}
	}
	if len(allowed) == 0 {
		log.Println("warning: ALLOWED_ORIGINS is not set — all cross-origin requests will be rejected")
	}

	isAllowed := func(origin string) bool {
		for _, o := range allowed {
			if o == origin {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && isAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
