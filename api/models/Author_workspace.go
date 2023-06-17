package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthorWorkspace struct {
	PersonID  uint64 `gorm:"primaryKey"`
	AddressID uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

