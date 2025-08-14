package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck responds with service health status.
func (*HTTPHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
