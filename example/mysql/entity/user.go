package entity

// User table struct.
type User struct {
	// ID    int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	// Name  string `gorm:"column:name;type:varchar(255);not null" mapstructure:"name"`
	// Email string `gorm:"column:email;type:varchar(255);not null" mapstructure:"email"`
	ID    int
	Name  string
	Email string
}

// TableName returns tablename.
func (i *User) TableName() string {
	return "user"
}
