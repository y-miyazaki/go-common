// Package db provides database transaction utilities.
package db

import (
	"fmt"

	"gorm.io/gorm"
)

// TransactionGorm automatically performs transactions (Begin/Commit/Rollback) between specified func...
// It starts a transaction, executes the provided function, and commits on success or rolls back on error.
func TransactionGorm(db *gorm.DB, f func(db *gorm.DB) error) error {
	// Track if transaction was committed successfully
	committed := false
	// Begin the transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin transaction: %w", tx.Error)
	}
	// Ensure rollback is called if transaction wasn't committed
	defer (func() {
		if !committed {
			tx.Rollback()
		}
	})()
	// Execute the provided function within the transaction
	if err := f(tx); err != nil {
		return fmt.Errorf("transaction function: %w", err)
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	// Mark as successfully committed
	committed = true
	return nil
}
