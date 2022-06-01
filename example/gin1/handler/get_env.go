package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetEnv handler
func (h *HTTPHandler) GetEnv(c *gin.Context) {
	password := os.Getenv("APP_DATABASE_MASTER_PASSWORD")
	addr := os.Getenv("REDIS_ADDR")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Hello!",
		"password": password,
		"addr":     addr,
	})
}
