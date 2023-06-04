package models

import "time"

type Workspace struct {
	ID          uint64    `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"column:name;size:255;not null;unique" json:"name"`
	Description string    `gorm:"column:description;size:255;not null;" json:"description"`
	Status      string    `gorm:"column:status;size:100;not null;default:'created'" json:"status"`
	Jobs        []Job     `gorm:"foreignKey:workspace_id"`
	Authors     []*Author `gorm:"many2many:author_workspace;"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
