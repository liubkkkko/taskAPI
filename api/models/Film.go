package models

import (
	"time"
)

type Film struct {
	ID         uint64  `gorm:"primary_key;auto_increment" json:"id"`
	Name       string  `gorm:"size:255;not null;unique" json:"name"`
	Genre      string  `gorm:"size:255;not null;" json:"genre"`
	DirectorId string  `gorm:"size:255;not null;" json:"directorId"`
	Rate       float32 `gorm:"size:255;not null;" json:"rate"`
	Year       uint16  `gorm:"size:255;not null;" json:"year"`
	Minutes uint16 `gorm:"size:255;not null;" json:"minutes"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// type Film struct {

// 	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
// 	Name  	  string    `gorm:"size:255;not null;unique" json:"nickname"`
// 	Genre     string    `gorm:"size:100;not null;unique" json:"email"`
// 	Director_id  string    `gorm:"size:100;not null;" json:"password"`
// 	Rate	string `gorm:"size:100;not null;" json:"password"`
// 	Year string `gorm:"size:100;not null;" json:"password"`
// 	Minutes string `gorm:"size:100;not null;" json:"password"`
// 	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
// 	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
// }
