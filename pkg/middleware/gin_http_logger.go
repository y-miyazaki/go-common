package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/infrastructure"
)

// GinHTTPLogger retrieves the request/response logs.
func GinHTTPLogger(
	logger *infrastructure.Logger,
	traceIDHeader string,
	clientIPHeader string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		fields := logrus.Fields{
			"host":      c.Request.Host,
			"duration":  duration.String(),
			"clientIP":  clientIP(c, clientIPHeader),
			"method":    c.Request.Method,
			"url":       c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"referer":   c.Request.Referer(),
			"userAgent": c.Request.UserAgent(),
		}
		if traceIDHeader != "" {
			fields[traceIDHeader] = c.Request.Header.Get(traceIDHeader)
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			l := logger
			var tmp interface{}
			tmp, _ = c.Get("error")
			if err, ok := tmp.(error); ok {
				l = logger.WithError(err)
			}
			messages, _ := c.Get("messages")
			l = l.WithField("messages", messages)
			l.WithFields(fields).Error()
		} else {
			logger.WithFields(fields).Info()
		}
	}
}

// clientIP gets ip address of client.
func clientIP(c *gin.Context, clientIPHeader string) string {
	if ip := c.Request.Header.Get(clientIPHeader); ip != "" {
		return ip
	}
	return c.ClientIP()
}
