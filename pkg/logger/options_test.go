package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestIsSensitiveKey(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"password", true},
		{"Password", true},
		{"api_key", true},
		{"session_token", true},
		{"username", false},
		{"email", false},
	}
	for _, tc := range tests {
		if got := isSensitiveKey(tc.key); got != tc.want {
			t.Fatalf("isSensitiveKey(%q) = %v; want %v", tc.key, got, tc.want)
		}
	}
}

func TestSanitizeFields_RedactsAndPreserves(t *testing.T) {
	fields := logrus.Fields{
		"password": "hunter2", // pragma: allowlist secret
		"user":     "alice",
	}
	got := SanitizeFields(fields, nil)
	if got["password"] != "[REDACTED]" {
		t.Fatalf("password was not redacted: %v", got["password"])
	}
	if got["user"] != "alice" {
		t.Fatalf("user was modified: %v", got["user"])
	}
}

func TestSanitizeFields_AllowSensitiveTrue(t *testing.T) {
	fields := logrus.Fields{"api_key": "abcd1234"} // pragma: allowlist secret
	cfg := &LoggerConfig{AllowSensitive: true}
	got := SanitizeFields(fields, cfg)
	if got["api_key"] != "abcd1234" { // pragma: allowlist secret
		t.Fatalf("expected api_key to be preserved when AllowSensitive=true, got %v", got["api_key"])
	}
}

func TestSanitizeFields_Nil(t *testing.T) {
	if got := SanitizeFields(nil, nil); got != nil {
		t.Fatalf("expected nil in, nil out; got %#v", got)
	}
}

func TestLogrusLogger_WithFieldRedaction_DefaultAndAllow(t *testing.T) {
	lr := logrus.New()
	l := NewLogger(lr)
	l2 := l.WithField("api_key", "secretval")
	if v, ok := l2.GetEntry().Data["api_key"]; !ok || v != "[REDACTED]" {
		t.Fatalf("expected api_key to be redacted by default, got %#v", v)
	}

	cfg := &LoggerConfig{AllowSensitive: true}
	l3 := NewLogger(lr, cfg)
	l4 := l3.WithField("api_key", "secretval")
	if v, ok := l4.GetEntry().Data["api_key"]; !ok || v != "secretval" {
		t.Fatalf("expected api_key to be preserved when AllowSensitive=true, got %#v", v)
	}
}
