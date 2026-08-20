//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package apigateway

import (
	"testing"
)

func TestGetSecureHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "x-frame-options", key: "X-Frame-Options", want: "DENY"},
		{name: "x-xss-protection", key: "X-XSS-Protection", want: "1; mode=block"},
		{name: "x-content-type-options", key: "X-Content-Type-Options", want: "nosniff"},
		{name: "strict-transport-security", key: "Strict-Transport-Security", want: "max-age=86400"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			headers := GetSecureHeaders()
			if headers == nil {
				t.Fatal("GetSecureHeaders() = nil, want non-nil map")
			}
			got, ok := headers[tt.key]
			if !ok {
				t.Fatalf("GetSecureHeaders()[%q] missing key", tt.key)
			}
			if got != tt.want {
				t.Fatalf("GetSecureHeaders()[%q] = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestGetCacheControlNoStoreHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
	}{
		{name: "cache-control no-store", want: "no-store"},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			headers := GetCacheControlNoStoreHeaders()
			if headers == nil {
				t.Fatal("GetCacheControlNoStoreHeaders() = nil, want non-nil map")
			}
			got, ok := headers["Cache-Control"]
			if !ok {
				t.Fatal("GetCacheControlNoStoreHeaders() missing Cache-Control key")
			}
			if got != tt.want {
				t.Fatalf("GetCacheControlNoStoreHeaders()[Cache-Control] = %q, want %q", got, tt.want)
			}
		})
	}
}
