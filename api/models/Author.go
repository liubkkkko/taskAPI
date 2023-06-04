package models

import (
	"time"
)

type Author struct {
	ID         uint64       `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Nickname   string       `gorm:"column:nickname;size:255;not null;" json:"nickname"`
	Email      string       `gorm:"column:email;size:100;not null;" json:"email"`
	Password   string       `gorm:"column:password;size:100;not null;" json:"password"`
	Role       string       `gorm:"column:role;size:100;not null;default:'user'" json:"role"`
	Jobs       []Job        `gorm:"foreignKey:author_id"`
	Workspaces []*Workspace `gorm:"many2many:author_workspace;"`
	CreatedAt  time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
