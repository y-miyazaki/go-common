// Package entity defines database entity models.
package entity

// User table struct.
type User struct {
	Name  string
	Email string
	ID    int
}

// TableName returns table name.
func (*User) TableName() string {
	return "user"
}
