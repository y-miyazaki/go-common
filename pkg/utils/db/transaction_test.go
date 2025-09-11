package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestTransactionGorm_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = TransactionGorm(gormDB, func(tx *gorm.DB) error {
		// Simulate successful operation
		return nil
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionGorm_FunctionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err = TransactionGorm(gormDB, func(tx *gorm.DB) error {
		// Simulate error in function
		return errors.New("function error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction function")
	assert.NoError(t, mock.ExpectationsWereMet())
}
