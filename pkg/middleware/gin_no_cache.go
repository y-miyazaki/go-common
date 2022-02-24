package middleware

import (
	"github.com/gin-gonic/gin"
)

// GinHTTPHeaderAccessControlAllowOrigin sets cache header.
// Cache-Control: no-store
// Pragma: no-cache
func GinHTTPHeaderNoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Next()
	}
}
