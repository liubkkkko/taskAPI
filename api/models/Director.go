package models

import (
	"time"
)

type Director struct {
	ID            uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name          string    `gorm:"size:255;not null;unique" json:"name"`
	Date_of_birth time.Time `json:"dateOfBirth"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
