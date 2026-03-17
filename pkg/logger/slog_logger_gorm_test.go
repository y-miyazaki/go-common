package logger

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ToDoSlog struct {
	ID   int
	Name string
}

func TestNewSlogGormLogger(t *testing.T) {
	loggerGorm := NewSlogLoggerGorm(&GormSlogSetting{
		Config: &SlogConfig{
			Level:  LevelInfo,
			Output: &bytes.Buffer{},
			Format: "json",
		},
		GormConfig: &GormConfig{
			SlowThreshold:             time.Duration(1000),
			IgnoreRecordNotFoundError: false,
			LogLevel:                  Info,
		},
	})

	mockdb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 loggerGorm,
	})

	todo := &ToDoSlog{Name: "testToDoSlog"}
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestNewSlogGormLoggerSlowQuery(t *testing.T) {
	loggerGorm := NewSlogLoggerGorm(&GormSlogSetting{
		Config: &SlogConfig{
			Level:  LevelInfo,
			Output: &bytes.Buffer{},
			Format: "json",
		},
		GormConfig: &GormConfig{
			SlowThreshold:             time.Duration(1),
			IgnoreRecordNotFoundError: false,
			LogLevel:                  Info,
		},
	})

	mockdb, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{Conn: mockdb}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 loggerGorm,
	})

	todo := &ToDoSlog{Name: "testToDoSlog"}
	db.Create(todo)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSlogGormLogger_LogMode(t *testing.T) {
	gormLogger := NewSlogLoggerGorm(&GormSlogSetting{
		Config: &SlogConfig{
			Level:  LevelInfo,
			Output: &bytes.Buffer{},
			Format: "json",
		},
		GormConfig: &GormConfig{LogLevel: Info},
	})

	result := gormLogger.LogMode(2)
	assert.Equal(t, gormLogger, result)
}

func TestSlogGormLogger_BasicMethods(t *testing.T) {
	gormLogger := NewSlogLoggerGorm(&GormSlogSetting{
		Config: &SlogConfig{
			Level:  LevelInfo,
			Output: &bytes.Buffer{},
			Format: "json",
		},
		GormConfig: &GormConfig{LogLevel: Info},
	})

	ctx := context.Background()
	gormLogger.Info(ctx, "Test info message: %s", "data")
	gormLogger.Warn(ctx, "Test warn message: %s", "data")
	gormLogger.Error(ctx, "Test error message: %s", "data")

	gormLogger.gormConfig.LogLevel = Silent
	gormLogger.Info(ctx, "This should not be logged")
	gormLogger.Warn(ctx, "This should not be logged")
	gormLogger.Error(ctx, "This should not be logged")
}

func TestSlogGormLogger_DefaultSettings(t *testing.T) {
	gormLogger := NewSlogLoggerGorm(nil)
	if gormLogger == nil {
		t.Fatal("expected logger to be created with default settings")
	}
	if gormLogger.gormConfig == nil {
		t.Fatal("expected default gorm config")
	}
	if gormLogger.gormConfig.LogLevel != Info {
		t.Fatalf("expected default LogLevel Info, got %v", gormLogger.gormConfig.LogLevel)
	}
}
