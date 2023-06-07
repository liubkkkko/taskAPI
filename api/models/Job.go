package models

import "time"

type Job struct {
	ID          uint64    `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Title       string    `gorm:"column:title;size:255;not null;" json:"title"`
	Content     string    `gorm:"column:content;mn:content;size:255;not null;" json:"content"`
	Status      string    `gorm:"column:status;size:255;not null;default:'created'" json:"status"`
	AuthorID    uint64    `gorm:"column:author_id;default:null" json:"author_id"`
	WorkspaceID uint64    `gorm:"column:workspace_id;default:null" json:"workspace_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}