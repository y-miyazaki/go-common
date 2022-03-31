package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHello handler
func (h *HTTPHandler) GetHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello!"})
}
