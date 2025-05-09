package logger

import (
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
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
