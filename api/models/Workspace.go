package models

import "time"

type Workspace struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name  string    `gorm:"size:255;not null;unique" json:"name"`
	Description     string    `gorm:"size:255;not null;" json:"description"`
	Status  string    `gorm:"size:100;not null;" json:"status"`
	Jobs []Job
	Authors []Author `gorm:"many2many:author_workspace;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
