package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/context"
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
		// get error
		l := logger
		if err, err2 := context.GetContextError(c); err2 == nil {
			l = l.WithError(err)
		}
		// get error message
		if messages, err := context.GetContextErrorMessage(c); err == nil {
			l = l.WithField("messages", messages)
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			l.WithFields(fields).Error()
		} else {
			l.WithFields(fields).Info()
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
