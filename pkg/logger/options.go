package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// LoggerConfig controls logger behavior related to sensitive data output.
type LoggerConfig struct {
	// AllowSensitive controls whether sensitive fields (password, token, etc.)
	// are allowed to be output in clear text. Default: false (do not allow).
	AllowSensitive bool
}

// defaultConfig returns a non-nil config with default values.
func defaultConfig(c *LoggerConfig) *LoggerConfig {
	if c == nil {
		return &LoggerConfig{AllowSensitive: false}
	}
	return c
}

// isSensitiveKey checks whether a key likely holds sensitive data.
// This is a conservative, case-insensitive substring match for common secret key names.
func isSensitiveKey(key string) bool {
	k := strings.ToLower(key)
	sensitiveSubstrings := []string{
		"password", "passwd", "pwd", "secret", "token", "access_token",
		"api_key", "apikey", "credential", "credentials", "auth", "private_key",
		"secret_key", "session_token", "clientip", "client_ip", "authorization", "set-cookie", "cookie", "x-forwarded-for",
	}
	for _, s := range sensitiveSubstrings {
		if strings.Contains(k, s) {
			return true
		}
	}
	return false
}

// SanitizeFields returns a copy of fields where sensitive keys are redacted
// according to the provided config. If cfg.AllowSensitive is true, fields
// are returned unchanged. Values that match sensitive patterns (e.g., Bearer
// tokens) are also redacted regardless of key name.
func SanitizeFields(fields logrus.Fields, cfg *LoggerConfig) logrus.Fields {
	cfgLocal := defaultConfig(cfg)
	if cfgLocal.AllowSensitive || fields == nil {
		return fields
	}
	out := logrus.Fields{}
	for k, v := range fields {
		if isSensitiveKey(k) {
			out[k] = "[REDACTED]"
			continue
		}
		if s, ok := v.(string); ok && isSensitiveValue(s) {
			out[k] = "[REDACTED]"
			continue
		}
		out[k] = v
	}
	return out
}

// isSensitiveValue checks whether a value looks like it contains sensitive
// data such as a Bearer token or a basic auth credential.
func isSensitiveValue(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(v, "bearer ") || strings.HasPrefix(v, "basic ") || strings.HasPrefix(v, "token ")
}

// SanitizeValue returns "[REDACTED]" if the key is sensitive or the value
// matches a sensitive pattern; otherwise it returns the value unchanged.
// Use this to sanitize individual header values before logging.
func SanitizeValue(key, value string) string {
	if isSensitiveKey(key) || isSensitiveValue(value) {
		return "[REDACTED]"
	}
	return value
}
