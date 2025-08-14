package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLoggerWithDifferentLevels(t *testing.T) {
	level := zap.NewAtomicLevel()
	zapConfig := zap.Config{
		Level: level,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "st",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	// Test Debug Level
	level.SetLevel(zap.DebugLevel)
	logger.Debug("Debug level test")
	logger.Info("Info level test")
	logger.Warn("Warn level test")
	logger.Error("Error level test")

	// Test Warn Level
	level.SetLevel(zap.WarnLevel)
	logger.Debug("This should not appear")
	logger.Info("This should not appear")
	logger.Warn("Warn level test")
	logger.Error("Error level test")

	// Test Error Level
	level.SetLevel(zap.ErrorLevel)
	logger.Debug("This should not appear")
	logger.Info("This should not appear")
	logger.Warn("This should not appear")
	logger.Error("Error level test")
}

func TestZapLoggerWithCustomFields(t *testing.T) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.InfoLevel)

	zapConfig := zap.Config{
		Level: level,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "st",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	logger.With(zap.String("customKey", "customValue")).Info("Testing custom fields")
	logger.With(zap.Int("userID", 12345)).Warn("Testing custom integer field")
	logger.With(zap.Bool("isActive", true)).Error("Testing custom boolean field")
}

func TestZapLoggerWithNilError(t *testing.T) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.InfoLevel)

	zapConfig := zap.Config{
		Level: level,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "st",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	var err error = nil
	logger.WithError(err).Error("Testing with nil error")
}

func TestZapLoggerWithPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic but did not occur")
		}
	}()

	level := zap.NewAtomicLevel()
	level.SetLevel(zap.PanicLevel)

	zapConfig := zap.Config{
		Level: level,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "st",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)
	logger.Panic("This should trigger a panic")
}
