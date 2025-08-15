// Package entity defines database entity models for PostgreSQL examples.
package entity

// User table struct.
type User struct {
	Name  string
	Email string
	ID    int
}

// TableName returns tablename.
func (*User) TableName() string {
	return "user"
}
