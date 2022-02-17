package infrastructure

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger struct.
type Logger struct {
	e *logrus.Entry
}

// NewLogger returns an instance of logger
func NewLogger(logger *logrus.Logger) *Logger {
	var l = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	l.SetFormatter(logger.Formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(logger.Out)
	// Only log the warning severity or above.
	l.SetLevel(logger.Level)
	return &Logger{
		e: l.WithFields(logrus.Fields{}),
	}
}

// WithField calls WithField function of logger entry.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		e: l.e.WithField(key, value),
	}
}

// WithFields calls WithField function of logger entry.
func (l *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		e: l.e.WithFields(fields),
	}
}

// WithError calls WithError function of logger entry.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		e: l.e.WithField("stacktrace", fmt.Sprintf("%+v", err)).WithError(err),
	}
}

// WithContext calls WithContext function of logger entry.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		e: l.e.WithContext(ctx),
	}
}

// Debugf outputs debug level log.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.e.Debugf(format, args...)
}

// Infof outputs info level log.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.e.Infof(format, args...)
}

// Printf outputs printf.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.e.Printf(format, args...)
}

// Warnf outputs warn level log.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.e.Warnf(format, args...)
}

// Warningf outputs warn level log.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.e.Warningf(format, args...)
}

// Errorf outputs error level log.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.e.Errorf(format, args...)
}

// Fatalf outputs fatal level log.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.e.Fatalf(format, args...)
}

// Panicf outputs panic log.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.e.Panicf(format, args...)
}
