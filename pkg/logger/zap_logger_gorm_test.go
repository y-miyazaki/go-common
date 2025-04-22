package logger

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ToDoZap struct {
	ID   int
	Name string
}

func TestNewZapGormLogger(t *testing.T) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	zapConfig := zap.Config{
		Level:             level,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
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
	loggerGorm := NewZapLoggerGorm(&GormZapSetting{
		Config: &zapConfig,
		GormConfig: &GormConfig{
			SlowThreshold:             time.Duration(1000),
			IgnoreRecordNotFoundError: false,
			LogLevel:                  Info,
		},
	})

	mockdb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 loggerGorm,
	})

	// Pointer to the ToDoZap struct to be used as an argument
	todo := &ToDoZap{
		Name: "testToDoZap",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewZapGormLoggerSilent(t *testing.T) {

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	zapConfig := zap.Config{
		Level:             level,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
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
	loggerGorm := NewZapLoggerGorm(&GormZapSetting{
		Config: &zapConfig,
		GormConfig: &GormConfig{
			SlowThreshold:             time.Duration(1000),
			IgnoreRecordNotFoundError: false,
			LogLevel:                  Info,
		},
	})

	mockdb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 loggerGorm,
	})

	// Pointer to the ToDoZap struct to be used as an argument
	todo := &ToDoZap{
		Name: "testToDoZap",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewZapGormLoggerSlowQuery(t *testing.T) {

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	zapConfig := zap.Config{
		Level:             level,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
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
	loggerGorm := NewZapLoggerGorm(&GormZapSetting{
		Config: &zapConfig,
		GormConfig: &GormConfig{
			SlowThreshold:             time.Duration(1),
			IgnoreRecordNotFoundError: false,
			LogLevel:                  Info,
		},
	})

	mockdb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 loggerGorm,
	})

	// Pointer to the ToDoZap struct to be used as an argument
	todo := &ToDoZap{
		Name: "testToDoZap",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
