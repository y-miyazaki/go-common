package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/y-miyazaki/go-common/pkg/middleware"
)

func main() {
	router := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(middleware.GinCors(
		&middleware.GinCorsConfig{
			AllowAllOrigins:  false,
			AllowOrigins:     []string{"https://foo.com"},
			AllowMethods:     []string{"PUT", "PATCH"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	router.Run()
}
