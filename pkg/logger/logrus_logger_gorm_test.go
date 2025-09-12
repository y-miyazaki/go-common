package logger

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ToDo struct {
	ID   int
	Name string
}

func TestNewGormLogger(t *testing.T) {

	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Info")

	loggerGorm := NewLoggerGorm(&GormSetting{
		Logger: logger,
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

	// Pointer to the ToDo struct to be used as an argument
	todo := &ToDo{
		Name: "testToDo",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewGormLoggerSilent(t *testing.T) {

	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Silent")

	loggerGorm := NewLoggerGorm(&GormSetting{
		Logger: logger,
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

	// Pointer to the ToDo struct to be used as an argument
	todo := &ToDo{
		Name: "testToDo",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewGormLoggerSlowQuery(t *testing.T) {

	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, _ = logrus.ParseLevel("Info")

	loggerGorm := NewLoggerGorm(&GormSetting{
		Logger: logger,
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

	// Pointer to the ToDo struct to be used as an argument
	todo := &ToDo{
		Name: "testToDo",
	}
	// mock.ExpectExec(".*").
	// 	// WithArgs(todo.Name).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGormLogger_LogMode(t *testing.T) {
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.InfoLevel

	gormLogger := NewLoggerGorm(&GormSetting{
		Logger: logger,
		GormConfig: &GormConfig{
			LogLevel: Info,
		},
	})

	// Test LogMode returns the same instance
	result := gormLogger.LogMode(2) // logger.Info level
	assert.Equal(t, gormLogger, result)
}

func TestGormLogger_Info(t *testing.T) {
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.InfoLevel

	gormLogger := NewLoggerGorm(&GormSetting{
		Logger: logger,
		GormConfig: &GormConfig{
			LogLevel: Info,
		},
	})

	ctx := context.Background()

	// Test Info with sufficient log level
	gormLogger.Info(ctx, "Test info message: %s", "data")

	// Test Info with insufficient log level
	gormLogger.gormConfig.LogLevel = Error
	gormLogger.Info(ctx, "This should not be logged")
}

func TestGormLogger_Warn(t *testing.T) {
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.WarnLevel

	gormLogger := NewLoggerGorm(&GormSetting{
		Logger: logger,
		GormConfig: &GormConfig{
			LogLevel: Warn,
		},
	})

	ctx := context.Background()

	// Test Warn with sufficient log level
	gormLogger.Warn(ctx, "Test warn message: %s", "data")

	// Test Warn with insufficient log level
	gormLogger.gormConfig.LogLevel = Error
	gormLogger.Warn(ctx, "This should not be logged")
}

func TestGormLogger_Error(t *testing.T) {
	logger := &logrus.Logger{}
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.ErrorLevel

	gormLogger := NewLoggerGorm(&GormSetting{
		Logger: logger,
		GormConfig: &GormConfig{
			LogLevel: Error,
		},
	})

	ctx := context.Background()

	// Test Error with sufficient log level
	gormLogger.Error(ctx, "Test error message: %s", "data")

	// Test Error with insufficient log level (Silent)
	gormLogger.gormConfig.LogLevel = Silent
	gormLogger.Error(ctx, "This should not be logged")
}
