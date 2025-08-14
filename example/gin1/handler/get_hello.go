package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SayHello responds with a greeting message.
func (*HTTPHandler) SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello!"})
}
