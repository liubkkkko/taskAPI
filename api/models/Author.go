package models

import "time"

type Author struct {
	ID         uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Nickname   string `gorm:"size:255;not null;unique" json:"nickname"`
	Email      string `gorm:"size:100;not null;unique" json:"email"`
	Password   string `gorm:"size:100;not null;" json:"password"`
	Role       string `gorm:"size:100;not null;" json:"role"`
	Jobs       []Job
	Workspaces []Workspace `gorm:"many2many:author_workspace;"`
	CreatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
