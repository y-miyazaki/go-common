package apigateway

// GetSecureHeaders gets secure headers.
func GetSecureHeaders() map[string]string {
	return map[string]string{
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"X-Content-Type-Options":    "nosniff",
		"Strict-Transport-Security": "max-age=86400",
	}
}

// GetCacheControlNoStoreHeaders gets Cache-Control header.
func GetCacheControlNoStoreHeaders() map[string]string {
	return map[string]string{
		"Cache-Control": "no-store",
	}
}
