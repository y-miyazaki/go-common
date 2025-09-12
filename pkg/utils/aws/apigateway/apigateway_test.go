package apigateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecureHeaders(t *testing.T) {
	headers := GetSecureHeaders()
	assert.NotNil(t, headers)
	assert.Equal(t, "DENY", headers["X-Frame-Options"])
	assert.Equal(t, "1; mode=block", headers["X-XSS-Protection"])
	assert.Equal(t, "nosniff", headers["X-Content-Type-Options"])
	assert.Equal(t, "max-age=86400", headers["Strict-Transport-Security"])
}

func TestGetCacheControlNoStoreHeaders(t *testing.T) {
	headers := GetCacheControlNoStoreHeaders()
	assert.NotNil(t, headers)
	assert.Equal(t, "no-store", headers["Cache-Control"])
}
