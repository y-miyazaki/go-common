package main

import (
	"fmt"
	"os"
	"time"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/example/gin1/handler"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"github.com/y-miyazaki/go-common/pkg/middleware"
)

func main() {
	l := &logrus.Logger{}
	// formatter
	l.Formatter = &logrus.JSONFormatter{}
	// out
	l.Out = os.Stdout
	// level
	level, err := logrus.ParseLevel("info")
	if err != nil {
		panic(fmt.Sprintf("level can't set %v", level))
	}
	l.Level = level
	logger := logger.NewLogger(l)
	h := handler.NewHTTPHandler(logger)

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
	router.Use(helmet.Default())
	router.Use(middleware.GinHTTPLogger(logger, "request-id", "test"))
	{
		router.GET("/healthcheck", h.GetHealthcheck)
		router.GET("/hello", h.GetHello)
		router.GET("/error_1", h.GetError1)
		router.GET("/error_2", h.GetError2)
		router.GET("/mysql", h.GetMySQL)
		router.GET("/postgres", h.GetPostgres)
		router.GET("/s3", h.GetS3)
	}
	err = router.Run()
	if err != nil {
		logger.WithError(err).Error("router.Run() error...")
	}
}
