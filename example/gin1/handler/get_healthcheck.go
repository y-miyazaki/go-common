package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealthcheck handler
func (h *HTTPHandler) GetHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
