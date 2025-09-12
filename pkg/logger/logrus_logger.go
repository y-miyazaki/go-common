package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Logger struct.
type Logger struct {
	Entry  *logrus.Entry
	Config *LoggerConfig
}

// NewLogger returns an instance of logger
func NewLogger(logger *logrus.Logger, cfg ...*LoggerConfig) *Logger {
	var l = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	l.SetFormatter(logger.Formatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	l.SetOutput(logger.Out)
	// Only log the warning severity or above.
	l.SetLevel(logger.Level)

	var config *LoggerConfig
	if len(cfg) > 0 {
		config = cfg[0]
	}

	return &Logger{
		Entry:  l.WithFields(logrus.Fields{}),
		Config: config,
	}
}

// GetEntry gets *logrus.Entry.
func (l *Logger) GetEntry() *logrus.Entry {
	return l.Entry
}

// WithField attaches a field to the logger.
func (l *Logger) WithField(key string, value any) *Logger {
	cfg := defaultConfig(l.Config)
	if !cfg.AllowSensitive && isSensitiveKey(key) {
		return &Logger{
			Entry:  l.Entry.WithField(key, "[REDACTED]"),
			Config: l.Config,
		}
	}
	return &Logger{
		Entry:  l.Entry.WithField(key, value),
		Config: l.Config,
	}
}

// WithFields calls WithField function of logger entry.
func (l *Logger) WithFields(fields logrus.Fields) *Logger {
	return &Logger{
		Entry:  l.Entry.WithFields(SanitizeFields(fields, l.Config)),
		Config: l.Config,
	}
}

// ... existing code ...

// WithError calls WithError function of logger entry.
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return &Logger{
		Entry:  l.Entry.WithError(err),
		Config: l.Config,
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
func (l *Logger) Debugf(format string, args ...any) {
	l.Entry.Debugf(format, args...)
}

// Infof outputs info level log.
func (l *Logger) Infof(format string, args ...any) {
	l.Entry.Infof(format, args...)
}

// Printf outputs printf.
func (l *Logger) Printf(format string, args ...any) {
	l.Entry.Printf(format, args...)
}

// Warnf outputs warn level log.
func (l *Logger) Warnf(format string, args ...any) {
	l.Entry.Warnf(format, args...)
}

// Warningf outputs warn level log.
func (l *Logger) Warningf(format string, args ...any) {
	l.Entry.Warningf(format, args...)
}

// Errorf outputs error level log.
func (l *Logger) Errorf(format string, args ...any) {
	l.Entry.Errorf(format, args...)
}

// Fatalf outputs fatal level log.
func (l *Logger) Fatalf(format string, args ...any) {
	l.Entry.Fatalf(format, args...)
}

// Panicf outputs panic log.
func (l *Logger) Panicf(format string, args ...any) {
	l.Entry.Panicf(format, args...)
}

// Debug outputs debug level log.
func (l *Logger) Debug(args ...any) {
	l.Entry.Debug(args...)
}

// Info outputs info level log.
func (l *Logger) Info(args ...any) {
	l.Entry.Info(args...)
}

// Print outputs printf.
func (l *Logger) Print(args ...any) {
	l.Entry.Print(args...)
}

// Warn outputs warn level log.
func (l *Logger) Warn(args ...any) {
	l.Entry.Warn(args...)
}

// Warning outputs warn level log.
func (l *Logger) Warning(args ...any) {
	l.Entry.Warning(args...)
}

// Error outputs error level log.
func (l *Logger) Error(args ...any) {
	l.Entry.Error(args...)
}

// Fatal outputs fatal level log.
func (l *Logger) Fatal(args ...any) {
	l.Entry.Fatal(args...)
}

// Panic outputs panic log.
func (l *Logger) Panic(args ...any) {
	l.Entry.Panic(args...)
}

// Debugln outputs debug level log.
func (l *Logger) Debugln(args ...any) {
	l.Entry.Debugln(args...)
}

// Infoln outputs info level log.
func (l *Logger) Infoln(args ...any) {
	l.Entry.Infoln(args...)
}

// Println outputs printf.
func (l *Logger) Println(args ...any) {
	l.Entry.Println(args...)
}

// Warnln outputs warn level log.
func (l *Logger) Warnln(args ...any) {
	l.Entry.Warnln(args...)
}

// Warningln outputs warn level log.
func (l *Logger) Warningln(args ...any) {
	l.Entry.Warningln(args...)
}

// Errorln outputs error level log.
func (l *Logger) Errorln(args ...any) {
	l.Entry.Errorln(args...)
}

// Fatalln outputs fatal level log.
func (l *Logger) Fatalln(args ...any) {
	l.Entry.Fatalln(args...)
}

// Panicln outputs panic log.
func (l *Logger) Panicln(args ...any) {
	l.Entry.Panicln(args...)
}
