package infrastructure

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Logger struct.
type Logger struct {
	Entry *logrus.Entry
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
		Entry: l.WithFields(logrus.Fields{}),
	}
}

// GetEntry gets *logrus.Entry.
func (l *Logger) GetEntry() *logrus.Entry {
	return l.Entry
}

// WithField calls WithField function of logger entry.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Entry: l.Entry.WithField(key, value),
	}
}

// WithFields calls WithField function of logger entry.
func (l *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		Entry: l.Entry.WithFields(fields),
	}
}

// WithError calls WithError function of logger entry.
func (l *Logger) WithError(err error) *Logger {
	if e, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
		return &Logger{
			Entry: l.Entry.WithField("stacktrace", fmt.Sprintf("%+v", e.StackTrace())).WithError(err),
		}
	}
	return &Logger{
		Entry: l.Entry.WithError(err),
	}
}

// WithContext calls WithContext function of logger entry.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Entry: l.Entry.WithContext(ctx),
	}
}

// WithContextValue calls WithField function of logger entry.
func (l *Logger) WithContextValue(key string) *Logger {
	return &Logger{
		Entry: l.Entry.WithField(key, l.Entry.Context.Value(key)),
	}
}

// Debugf outputs debug level log.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Entry.Debugf(format, args...)
}

// Infof outputs info level log.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Entry.Infof(format, args...)
}

// Printf outputs printf.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Entry.Printf(format, args...)
}

// Warnf outputs warn level log.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

// Warningf outputs warn level log.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Entry.Warningf(format, args...)
}

// Errorf outputs error level log.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}

// Fatalf outputs fatal level log.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Entry.Fatalf(format, args...)
}

// Panicf outputs panic log.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Entry.Panicf(format, args...)
}

// Debug outputs debug level log.
func (l *Logger) Debug(args ...interface{}) {
	l.Entry.Debug(args...)
}

// Info outputs info level log.
func (l *Logger) Info(args ...interface{}) {
	l.Entry.Info(args...)
}

// Print outputs printf.
func (l *Logger) Print(args ...interface{}) {
	l.Entry.Print(args...)
}

// Warn outputs warn level log.
func (l *Logger) Warn(args ...interface{}) {
	l.Entry.Warn(args...)
}

// Warning outputs warn level log.
func (l *Logger) Warning(args ...interface{}) {
	l.Entry.Warning(args...)
}

// Error outputs error level log.
func (l *Logger) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

// Fatal outputs fatal level log.
func (l *Logger) Fatal(args ...interface{}) {
	l.Entry.Fatal(args...)
}

// Panic outputs panic log.
func (l *Logger) Panic(args ...interface{}) {
	l.Entry.Panic(args...)
}

// Debugln outputs debug level log.
func (l *Logger) Debugln(args ...interface{}) {
	l.Entry.Debugln(args...)
}

// Infoln outputs info level log.
func (l *Logger) Infoln(args ...interface{}) {
	l.Entry.Infoln(args...)
}

// Println outputs printf.
func (l *Logger) Println(args ...interface{}) {
	l.Entry.Println(args...)
}

// Warnln outputs warn level log.
func (l *Logger) Warnln(args ...interface{}) {
	l.Entry.Warnln(args...)
}

// Warningln outputs warn level log.
func (l *Logger) Warningln(args ...interface{}) {
	l.Entry.Warningln(args...)
}

// Errorln outputs error level log.
func (l *Logger) Errorln(args ...interface{}) {
	l.Entry.Errorln(args...)
}

// Fatalln outputs fatal level log.
func (l *Logger) Fatalln(args ...interface{}) {
	l.Entry.Fatalln(args...)
}

// Panicln outputs panic log.
func (l *Logger) Panicln(args ...interface{}) {
	l.Entry.Panicln(args...)
}
