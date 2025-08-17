// Package model provides base database model structures using GORM.
package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// BaseModel struct
// nolint:tagliatelle
type BaseModel struct {
	CreatedAt time.Time `mapstructure:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `mapstructure:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
}

// BaseModelWithDeleted struct
// nolint:tagliatelle
type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt soft_delete.DeletedAt `mapstructure:"deleted_at" gorm:"column:deleted_at;type:int(11);default:0"`
}
