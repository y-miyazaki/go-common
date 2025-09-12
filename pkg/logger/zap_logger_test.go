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

func TestZapLogger_PrintfStyleFunctions(t *testing.T) {
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	// Test printf-style functions
	logger.Debugf("Debugf test: %s", "debug")
	logger.Infof("Infof test: %s", "info")
	logger.Warnf("Warnf test: %s", "warn")
	logger.Errorf("Errorf test: %s", "error")

	// Note: Fatalf and Panicf would exit the program or panic, so we don't test them
	// logger.Fatalf("Fatalf test: %s", "fatal")
	// logger.Panicf("Panicf test: %s", "panic")

	// Test DPanicf (only panics in development mode)
	logger.DPanicf("DPanicf test: %s", "dpanic")
}

func TestZapLogger_PrintlnStyleFunctions(t *testing.T) {
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	// Test println-style functions
	logger.Debugln("Debugln test")
	logger.Infoln("Infoln test")
	logger.Warnln("Warnln test")
	logger.Errorln("Errorln test")

	// Note: Fatalln and Panicln would exit the program or panic, so we don't test them
	// logger.Fatalln("Fatalln test")
	// logger.Panicln("Panicln test")

	// Test DPanicln (only panics in development mode)
	logger.DPanicln("DPanicln test")
}

func TestZapLogger_FatalAndDPanic(t *testing.T) {
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false, // Set to false so DPanic doesn't panic
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := NewZapLogger(&zapConfig)

	// Test Fatal - this would normally exit, but in test we can call it
	// Note: In real usage, Fatal would exit the program
	// For testing purposes, we just ensure the function exists and can be called
	// logger.Fatal("Fatal test")

	// Test DPanic - only panics in development mode
	logger.DPanic("DPanic test")
}
