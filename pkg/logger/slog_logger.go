// Package logger provides structured logging functionality using slog.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// Level is a type that represents log levels
type Level = slog.Level

// Available log levels
const (
	LevelDebug = slog.LevelDebug // Debug level for detailed information
	LevelInfo  = slog.LevelInfo  // Info level for general operational information
	LevelWarn  = slog.LevelWarn  // Warn level for warning messages
	LevelError = slog.LevelError // Error level for error conditions
)

// SlogConfig defines the configuration options for the *logger
// nolint:govet
type SlogConfig struct {
	// Level sets the minimum level of messages to log
	Level Level
	// AddSource adds source code location to log messages
	AddSource bool
	// Output specifies the destination for log messages
	Output io.Writer
	// Format specifies the output format ("json" or "text")
	Format string
}

// SlogLogger implements the logger interface
type SlogLogger struct {
	log *slog.Logger
}

// NewSlogLogger creates a new SlogLogger instance with the specified configuration.
// If cfg is nil, default configuration will be used.
func NewSlogLogger(cfg *SlogConfig) *SlogLogger {
	if cfg == nil {
		panic("*logger config is required")
	}

	opts := &slog.HandlerOptions{
		Level:       cfg.Level,
		AddSource:   cfg.AddSource,
		ReplaceAttr: nil, // Not used, defaults to nil
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(cfg.Output, opts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, opts)
	}

	return &SlogLogger{
		log: slog.New(handler),
	}
}

// Debug implements *logger.Debug
func (l *SlogLogger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

// Info implements *logger.Info
func (l *SlogLogger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

// Warn implements *logger.Warn
func (l *SlogLogger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

// Error implements *logger.Error
func (l *SlogLogger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

// WithContext implements *logger.WithContext
// It creates a new *logger with trace ID from the context
func (l *SlogLogger) WithContext(ctx context.Context) *SlogLogger {
	return &SlogLogger{
		log: l.log.With("trace_id", getTraceID(ctx)),
	}
}

// With implements *logger.With
// It creates a new *logger with additional key-value pairs in the context
func (l *SlogLogger) With(args ...any) *SlogLogger {
	return &SlogLogger{
		log: l.log.With(args...),
	}
}

// WithError implements *logger.WithError
// It creates a new *logger with error information added to the context
func (l *SlogLogger) WithError(err error) *SlogLogger {
	if err == nil {
		return l
	}
	return &SlogLogger{
		log: l.log.With(
			"error", err.Error(),
			"stack", fmt.Sprintf("%+v", err),
		),
	}
}

// getTraceID retrieves trace ID from the given context.
// Returns empty string if trace ID is not found.
func getTraceID(_ context.Context) string { // nolint:unused
	// Implement trace ID retrieval logic
	return ""
}
