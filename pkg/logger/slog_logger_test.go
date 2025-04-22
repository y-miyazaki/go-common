package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// TestNewLogger tests logger initialization
func TestNewSlogLogger(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *SlogConfig
		wantPanic bool
	}{
		{
			name:      "should panic with nil config",
			cfg:       nil,
			wantPanic: true,
		},
		{
			name: "should work with JSON format",
			cfg: &SlogConfig{
				Level:     LevelInfo,
				AddSource: true,
				Output:    os.Stdout,
				Format:    "json",
			},
			wantPanic: false,
		},
		{
			name: "should work with text format",
			cfg: &SlogConfig{
				Level:     LevelInfo,
				AddSource: true,
				Output:    os.Stdout,
				Format:    "text",
			},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic did not occur")
					}
				}()
			}
			log := NewSlogLogger(tt.cfg)
			if !tt.wantPanic && log == nil {
				t.Error("failed to create logger")
			}
		})
	}
}

// TestLogger_LogLevels tests log level outputs
func TestLogger_LogLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelDebug,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	tests := []struct {
		name      string
		logFunc   func()
		wantLevel string
		wantMsg   string
	}{
		{
			name: "should output debug level log",
			logFunc: func() {
				log.Debug("debug message", "key", "value")
			},
			wantLevel: "DEBUG",
			wantMsg:   "debug message",
		},
		{
			name: "should output info level log",
			logFunc: func() {
				log.Info("info message", "key", "value")
			},
			wantLevel: "INFO",
			wantMsg:   "info message",
		},
		{
			name: "should output warn level log",
			logFunc: func() {
				log.Warn("warning message", "key", "value")
			},
			wantLevel: "WARN",
			wantMsg:   "warning message",
		},
		{
			name: "should output error level log",
			logFunc: func() {
				log.Error("error message", "key", "value")
			},
			wantLevel: "ERROR",
			wantMsg:   "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			var result map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("failed to parse log output: %v", err)
			}

			// Verify results
			if level := result["level"]; level != tt.wantLevel {
				t.Errorf("log level mismatch: got=%v, want=%v", level, tt.wantLevel)
			}
			if msg := result["msg"]; msg != tt.wantMsg {
				t.Errorf("message mismatch: got=%v, want=%v", msg, tt.wantMsg)
			}
		})
	}
}

// TestLogger_WithContext tests logger with context
func TestLogger_WithContext(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelInfo,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	ctx := context.Background()
	logWithCtx := log.WithContext(ctx)
	logWithCtx.Info("test message")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if _, exists := result["trace_id"]; !exists {
		t.Error("trace_id field not found")
	}
}

// TestLogger_WithError tests logger with error information
func TestLogger_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelInfo,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	testErr := errors.New("test error")
	logWithErr := log.WithError(testErr)
	logWithErr.Error("operation failed")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if errMsg := result["error"]; errMsg != testErr.Error() {
		t.Errorf("error message mismatch: got=%v, want=%v", errMsg, testErr.Error())
	}
	if stack := result["stack"]; !strings.Contains(stack.(string), "test error") {
		t.Errorf("stack trace does not contain error message: %v", stack)
	}
}

// TestLogger_With tests logger with additional key-value pairs
func TestLogger_With(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelInfo,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	logWith := log.With("key1", "value1", "key2", "value2")
	logWith.Info("test message")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if result["key1"] != "value1" || result["key2"] != "value2" {
		t.Errorf("key-value pairs not found in log: got=%v", result)
	}
}

// TestLogger_WithError_NilError tests logger behavior when WithError is called with nil
func TestLogger_WithError_NilError(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelInfo,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	logWithErr := log.WithError(nil)
	logWithErr.Error("operation failed")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if _, exists := result["error"]; exists {
		t.Error("error field should not exist when nil error is passed")
	}
}

// TestLogger_WithContext_NoTraceID tests logger behavior when context has no trace ID
func TestLogger_WithContext_NoTraceID(t *testing.T) {
	buf := &bytes.Buffer{}
	multiWriter := io.MultiWriter(buf, os.Stdout)

	log := NewSlogLogger(&SlogConfig{
		Level:     LevelInfo,
		AddSource: false,
		Output:    multiWriter,
		Format:    "json",
	})

	ctx := context.Background()
	logWithCtx := log.WithContext(ctx)
	logWithCtx.Info("test message")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if traceID, exists := result["trace_id"]; !exists || traceID != "" {
		t.Errorf("trace_id mismatch: got=%v, want empty string", traceID)
	}
}
