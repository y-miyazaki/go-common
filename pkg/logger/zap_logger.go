package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger struct.
type ZapLogger struct {
	Logger *zap.Logger
}

// NewZapLogger returns an instance of logger
func NewZapLogger(config *zap.Config) *ZapLogger {
	// Encoding
	if config.Encoding == "" {
		config.Encoding = "json"
	}

	// OutputPaths
	if len(config.OutputPaths) == 0 {
		config.OutputPaths = []string{"stdout"}
	}

	// ErrorOutputPaths
	if len(config.ErrorOutputPaths) == 0 {
		config.ErrorOutputPaths = []string{"stderr"}
	}

	logger, err := config.Build()
	if err != nil {
		panic("can't create logger from zap")
	}
	return &ZapLogger{Logger: logger}
}

// With calls WithField function of logger.
func (l *ZapLogger) With(fields ...zapcore.Field) *ZapLogger {
	return &ZapLogger{
		Logger: l.Logger.With(fields...),
	}
}

// WithError calls WithError function of logger.
func (l *ZapLogger) WithError(err error) *ZapLogger {
	if err == nil {
		return l
	}
	return &ZapLogger{
		Logger: l.Logger.With(zap.Error(err)),
	}
}

// Debugf outputs debug level log.
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Sugar().Debugf(format, args...)
}

// Infof outputs info level log.
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.Logger.Sugar().Infof(format, args...)
}

// Warnf outputs warn level log.
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Sugar().Warnf(format, args...)
}

// Errorf outputs error level log.
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Sugar().Errorf(format, args...)
}

// Fatalf outputs fatal level log.
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Sugar().Fatalf(format, args...)
}

// Panicf outputs panic log.
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.Logger.Sugar().Panicf(format, args...)
}

// DPanicf outputs panic log.
func (l *ZapLogger) DPanicf(format string, args ...interface{}) {
	l.Logger.Sugar().DPanicf(format, args...)
}

// Debug outputs debug level log.
func (l *ZapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info outputs info level log.
func (l *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn outputs warn level log.
func (l *ZapLogger) Warn(msg string, fields ...zapcore.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error outputs error level log.
func (l *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal outputs fatal level log.
func (l *ZapLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Panic outputs panic log.
func (l *ZapLogger) Panic(msg string, fields ...zapcore.Field) {
	l.Logger.Panic(msg, fields...)
}

// DPanic outputs panic log.
func (l *ZapLogger) DPanic(msg string, fields ...zapcore.Field) {
	l.Logger.DPanic(msg, fields...)
}

// Debugln outputs debug level log.
func (l *ZapLogger) Debugln(args ...interface{}) {
	l.Logger.Sugar().Debugln(args...)
}

// Infoln outputs info level log.
func (l *ZapLogger) Infoln(args ...interface{}) {
	l.Logger.Sugar().Infoln(args...)
}

// Warnln outputs warn level log.
func (l *ZapLogger) Warnln(args ...interface{}) {
	l.Logger.Sugar().Warnln(args...)
}

// Errorln outputs error level log.
func (l *ZapLogger) Errorln(args ...interface{}) {
	l.Logger.Sugar().Errorln(args...)
}

// Fatalln outputs fatal level log.
func (l *ZapLogger) Fatalln(args ...interface{}) {
	l.Logger.Sugar().Fatalln(args...)
}

// Panicln outputs panic log.
func (l *ZapLogger) Panicln(args ...interface{}) {
	l.Logger.Sugar().Panicln(args...)
}

// DPanicln outputs panic log.
func (l *ZapLogger) DPanicln(args ...interface{}) {
	l.Logger.Sugar().DPanicln(args...)
}
