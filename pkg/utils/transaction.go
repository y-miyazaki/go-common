package utils

import "gorm.io/gorm"

// TransactionGorm automatically performs transactions (Begin/Commit/Rollback) between specified func...
func TransactionGorm(db *gorm.DB, f func(db *gorm.DB) error) error {
	committed := false
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer (func() {
		if !committed {
			tx.Rollback()
		}
	})()
	if err := f(tx); err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true
	return nil
}
