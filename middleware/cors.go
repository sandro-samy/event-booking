package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORS allows the configured CORS_ORIGIN to call the API. When it isn't set
// (the default for local development) it reflects back whatever Origin the
// browser sent, so the frontend dev server works regardless of which port it
// picks, without needing any env configuration.
func CORS(c *gin.Context) {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = c.GetHeader("Origin")
	}
	if origin == "" {
		c.Next()
		return
	}

	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Vary", "Origin")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
