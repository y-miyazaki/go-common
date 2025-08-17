// Package model provides base database model structures using GORM.
package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// BaseModel struct
// nolint:tagliatelle
type BaseModel struct {
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" mapstructure:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" mapstructure:"updated_at"`
}

// BaseModelWithDeleted struct
// nolint:tagliatelle
type BaseModelWithDeleted struct {
	BaseModel
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int(11);default:0" mapstructure:"deleted_at"`
}
