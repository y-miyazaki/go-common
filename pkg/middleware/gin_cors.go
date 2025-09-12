// Package middleware provides HTTP middleware components for Gin web framework.
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	decimal int = 10
)

type converter func(string) string // nolint:unused // type definition for potential future use

// GinCorsConfig sets configurations.
type GinCorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           time.Duration
	AllowAllOrigins  bool
	AllowCredentials bool
}

// GinCors sets Access-Control-XXXXX header.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
func GinCors(
	cs *GinCorsConfig,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			// request is not a CORS request.
			c.Next()
			return
		}

		// validate origin
		if !cs.validateOrigin(origin) {
			c.AbortWithStatus(http.StatusForbidden)
			c.Next()
			return
		}

		// Preflight headers
		if c.Request.Method == http.MethodOptions {
			// Access-Control-Allow-Methods
			if len(cs.AllowMethods) > 0 {
				allowMethods := convert(normalize(cs.AllowMethods), strings.ToUpper)
				c.Header("Access-Control-Allow-Methods", strings.Join(allowMethods, ","))
			}
			// Access-Control-Allow-Headers
			if len(cs.AllowHeaders) > 0 {
				allowHeaders := convert(normalize(cs.AllowHeaders), http.CanonicalHeaderKey)
				c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
			}
			// Access-Control-Max-Age
			if cs.MaxAge > time.Duration(0) {
				value := strconv.FormatInt(int64(cs.MaxAge/time.Second), decimal)
				c.Header("Access-Control-Max-Age", value)
			}
		}
		// Access-Control-Allow-Origin
		if cs.AllowAllOrigins {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Add("Vary", "Origin")
			c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
			c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")
		}

		// Access-Control-Allow-Credentials
		if cs.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Access-Control-Expose-Headers
		if len(cs.ExposeHeaders) > 0 {
			exposeHeaders := convert(normalize(cs.ExposeHeaders), http.CanonicalHeaderKey)
			c.Header("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ","))
		}

		// Check OPTIONS
		if c.Request.Method == http.MethodOptions {
			if cs.validateMethodOptions() {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}
		c.Next()
	}
}

// DefaultConfig returns a generic default configuration mapped to localhost.
func DefaultConfig() *GinCorsConfig {
	const defaultMaxAgeSecondsStandard = 86400
	return &GinCorsConfig{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           defaultMaxAgeSecondsStandard * time.Second,
	}
}

func normalize(values []string) []string {
	if values == nil {
		return nil
	}
	distinctMap := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)
		if _, seen := distinctMap[value]; !seen {
			normalized = append(normalized, value)
			distinctMap[value] = true
		}
	}
	return normalized
}

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

func (cs *GinCorsConfig) validateOrigin(origin string) bool {
	if cs.AllowAllOrigins {
		return true
	}
	for _, v := range cs.AllowOrigins {
		if v == origin {
			return true
		}
	}
	return false
}
func (cs *GinCorsConfig) validateMethodOptions() bool {
	for _, v := range cs.AllowMethods {
		if v == "OPTIONS" {
			return true
		}
	}
	return false
}
