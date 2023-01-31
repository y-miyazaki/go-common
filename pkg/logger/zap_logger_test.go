package logger

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TestZap struct {
	name string
	sex  string
}

func TestNewZapLogger(t *testing.T) {
	testZap := TestZap{
		name: "Joe",
		sex:  "man",
	}
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
	}
	logger := NewZapLogger(zapConfig)
	logger.Debug("test", zap.String("Key", "String"))
	logger.Info("test")
	logger.Error("test", zap.String("Key", "String"))
	logger.Warn("test")

	zapConfig2 := zap.Config{
		Level:             level,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
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

	logger2 := NewZapLogger(zapConfig2)
	logger2.Debug("test", zap.String("Key", "String"))
	logger2.Info("test")
	logger2.Error("test", zap.String("Key", "String"), zap.Error(nil))
	logger2.Warn("test")
	logger2.Debugln("test")
	logger2.Debugln("test")
	e1 := errors.WithStack(errors.New("test1"))
	e2 := errors.New("test2")
	e3 := fmt.Errorf("test3 %s", e1)
	logger2.WithError(e1).Error("test1")
	logger2.WithError(e2).Error("test2")
	logger2.WithError(e3).Error("test3")
	logger2.With(zap.String("key", "test")).With(zap.Reflect("any", testZap)).Error("with test")
	logger2.Error("with test", zap.Any("any", &testZap))
	logger2.Info("info test", zap.String("name", testZap.name), zap.String("sex", testZap.sex))

	logger2.Debugf("Debugf")
	logger2.Infof("Infof")
	logger2.Warnf("Warnf")
	logger2.Errorf("Errorf")
	// logger2.Fatalf("Fatalf")
	// logger2.Panicf("Panicf")
	logger2.DPanicf("DPanicf")
	logger2.Debug("Debug")
	logger2.Info("Info")
	logger2.Warn("Warn")
	logger2.Error("Error")
	// logger2.Fatal("Fatal")
	// logger2.Panic("Panic")
	logger2.DPanic("DPanic")
	logger2.Debugln("Debugln")
	logger2.Infoln("Infoln")
	logger2.Warnln("Warnln")
	logger2.Errorln("Errorln")
	// logger2.Fatalln("Fatalln")
	// logger2.Panicln("Panicln")
	logger2.DPanicln("DPanicln")
}
