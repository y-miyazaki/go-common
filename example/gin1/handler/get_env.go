// Package handler provides HTTP request handlers for the Gin web framework.
package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// HandleEnv retrieves environment variables and returns them as JSON response.
func (*HTTPHandler) HandleEnv(c *gin.Context) {
	password := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	addr := os.Getenv("REDIS_ADDR")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Hello!",
		"password": password, // pragma: allowlist secret
		"addr":     addr,
	})
}
