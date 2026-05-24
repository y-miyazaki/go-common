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
	var cfg zap.Config
	if config == nil {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = *config
	}

	// Encoding
	if cfg.Encoding == "" {
		cfg.Encoding = "json"
	}

	// OutputPaths
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stdout"}
	}

	// ErrorOutputPaths
	if len(cfg.ErrorOutputPaths) == 0 {
		cfg.ErrorOutputPaths = []string{"stderr"}
	}

	logger, err := cfg.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	return &ZapLogger{Logger: logger}
}

// DPanic outputs panic log.
func (l *ZapLogger) DPanic(msg string, fields ...zapcore.Field) {
	l.Logger.DPanic(msg, fields...)
}

// DPanicf outputs panic log.
func (l *ZapLogger) DPanicf(format string, args ...any) {
	l.Logger.Sugar().DPanicf(format, args...)
}

// DPanicln outputs panic log.
func (l *ZapLogger) DPanicln(args ...any) {
	l.Logger.Sugar().DPanicln(args...)
}

// Debug outputs debug level log.
func (l *ZapLogger) Debug(msg string, fields ...zapcore.Field) {
	l.Logger.Debug(msg, fields...)
}

// Debugf outputs debug level log.
func (l *ZapLogger) Debugf(format string, args ...any) {
	l.Logger.Sugar().Debugf(format, args...)
}

// Debugln outputs debug level log.
func (l *ZapLogger) Debugln(args ...any) {
	l.Logger.Sugar().Debugln(args...)
}

// Error outputs error level log.
func (l *ZapLogger) Error(msg string, fields ...zapcore.Field) {
	l.Logger.Error(msg, fields...)
}

// Errorf outputs error level log.
func (l *ZapLogger) Errorf(format string, args ...any) {
	l.Logger.Sugar().Errorf(format, args...)
}

// Errorln outputs error level log.
func (l *ZapLogger) Errorln(args ...any) {
	l.Logger.Sugar().Errorln(args...)
}

// Fatal outputs fatal level log.
func (l *ZapLogger) Fatal(msg string, fields ...zapcore.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Fatalf outputs fatal level log.
func (l *ZapLogger) Fatalf(format string, args ...any) {
	l.Logger.Sugar().Fatalf(format, args...)
}

// Fatalln outputs fatal level log.
func (l *ZapLogger) Fatalln(args ...any) {
	l.Logger.Sugar().Fatalln(args...)
}

// Info outputs info level log.
func (l *ZapLogger) Info(msg string, fields ...zapcore.Field) {
	l.Logger.Info(msg, fields...)
}

// Infof outputs info level log.
func (l *ZapLogger) Infof(format string, args ...any) {
	l.Logger.Sugar().Infof(format, args...)
}

// Infoln outputs info level log.
func (l *ZapLogger) Infoln(args ...any) {
	l.Logger.Sugar().Infoln(args...)
}

// Panic outputs panic log.
func (l *ZapLogger) Panic(msg string, fields ...zapcore.Field) {
	l.Logger.Panic(msg, fields...)
}

// Panicf outputs panic log.
func (l *ZapLogger) Panicf(format string, args ...any) {
	l.Logger.Sugar().Panicf(format, args...)
}

// Panicln outputs panic log.
func (l *ZapLogger) Panicln(args ...any) {
	l.Logger.Sugar().Panicln(args...)
}

// Warn outputs warn level log.
func (l *ZapLogger) Warn(msg string, fields ...zapcore.Field) {
	l.Logger.Warn(msg, fields...)
}

// Warnf outputs warn level log.
func (l *ZapLogger) Warnf(format string, args ...any) {
	l.Logger.Sugar().Warnf(format, args...)
}

// Warnln outputs warn level log.
func (l *ZapLogger) Warnln(args ...any) {
	l.Logger.Sugar().Warnln(args...)
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
