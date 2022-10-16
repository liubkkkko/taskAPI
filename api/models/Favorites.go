package models

import (
	"time"
)

type Favorites struct {

	User_id            uint64    `gorm:"primary_key;auto_increment" json:"user_id"`
	Film_id          uint64    `gorm:"size:255;not null;unique" json:"film_id"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Wishlist struct {

	User_id            uint64    `gorm:"primary_key;auto_increment" json:"user_id"`
	Film_id          uint64    `gorm:"size:255;not null;unique" json:"film_id"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

