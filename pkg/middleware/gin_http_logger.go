// Package middleware provides HTTP middleware utilities for the Gin web framework.
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/y-miyazaki/go-common/pkg/context"
	"github.com/y-miyazaki/go-common/pkg/logger"
	"go.uber.org/zap"
)

// GinHTTPLogger retrieves the request/response logs.
// It logs HTTP requests and responses using logrus logger with configurable headers.
func GinHTTPLogger(l *logger.Logger, traceIDHeader, clientIPHeader string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time for duration calculation
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
		loggerWithContext := l
		if err, err2 := context.GetGinContextError(c); err2 == nil {
			loggerWithContext = loggerWithContext.WithError(err)
		}
		// get error message
		if messages, err := context.GetGinContextErrorMessage(c); err == nil {
			loggerWithContext = loggerWithContext.WithField("messages", messages)
		}
		if c.Writer.Status() >= http.StatusInternalServerError {
			loggerWithContext.WithFields(fields).Error()
		} else {
			loggerWithContext.WithFields(fields).Info()
		}
	}
}

// GinHTTPZapLogger retrieves the request/response logs.
func GinHTTPZapLogger(
	l *logger.ZapLogger,
	traceIDHeader, clientIPHeader string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		loggerWithContext := l.With(
			zap.String("host", c.Request.Host),
			zap.String("duration", duration.String()),
			zap.String("clientIP", clientIP(c, clientIPHeader)),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.RequestURI),
			zap.Int("status", c.Writer.Status()),
			zap.String("referer", c.Request.Referer()),
			zap.String("userAgent", c.Request.UserAgent()),
		)

		if traceIDHeader != "" {
			loggerWithContext = loggerWithContext.With(zap.String(traceIDHeader, c.Request.Header.Get(traceIDHeader)))
		}
		// get error
		if err, err2 := context.GetGinContextError(c); err2 == nil {
			loggerWithContext = loggerWithContext.WithError(err)
		}
		// get error message
		if messages, err := context.GetGinContextErrorMessage(c); err == nil {
			loggerWithContext = loggerWithContext.With(zap.String("messages", messages))
		}
		if c.Writer.Status() >= http.StatusInternalServerError {
			loggerWithContext.Error("")
		} else {
			loggerWithContext.Info("")
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
