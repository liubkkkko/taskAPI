package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

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

func (j *Job) Prepare() {
	j.ID = 0
	j.Title = html.EscapeString(strings.TrimSpace(j.Title))
	j.Content = html.EscapeString(strings.TrimSpace(j.Content))
	j.Status = html.EscapeString(strings.TrimSpace(j.Status))
	j.AuthorID = 0
	j.WorkspaceID = 0
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

func (j *Job) Validate() error {

	if j.Title == "" {
		return echo.NewHTTPError(418, "required Title")
	}
	if j.Content == "" {
		return echo.NewHTTPError(418, "required Content")
	}
	if j.AuthorID < 1 {
		return echo.NewHTTPError(418, "required Author")
	}
	if j.WorkspaceID < 1 {
		return echo.NewHTTPError(418, "required Workspace")
	}
	if j.Status == "" {
		return echo.NewHTTPError(418, "required Status")
	}
	return nil
}

func (j *Job) SaveJob(db *gorm.DB) (*Job, error) {
	err := db.Debug().Model(&Post{}).Create(&j).Error
	if err != nil {
		return &Job{}, err
	}
	// if j.ID != 0 {
	// 	err = db.Debug().Model(&Task{}).Where("id = ?", j.AuthorID).Take(&j.Author).Error
	// 	if err != nil {
	// 		return &Job{}, err
	// 	}
	// }
	return j, nil
}

func (j *Job) FindAllJob(db *gorm.DB) (*[]Job, error) {
	jobs := []Job{}
	err := db.Debug().Model(&Post{}).Limit(100).Find(&j).Error
	if err != nil {
		return &[]Job{}, err
	}
	// if len(jobs) > 0 {
	// 	for i := range jobs {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", tasks[i].AuthorID).Take(&tasks[i].Author).Error
	// 		if err != nil {
	// 			return &[]Task{}, err
	// 		}
	// 	}
	// }
	return &jobs, nil
}

func (j *Job) FindJobByID(db *gorm.DB, jid uint64) (*Job, error) {
	err := db.Debug().Model(&Task{}).Where("id = ?", jid).Take(&j).Error
	if err != nil {
		return &Job{}, err
	}
	// if t.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
	// 	if err != nil {
	// 		return &Task{}, err
	// 	}
	// }
	return j, nil
}

func (j *Job) UpdateAJob(db *gorm.DB) (*Job, error) {
	err := db.Debug().Model(&Job{}).Where("id = ?", j.ID).Updates(Job{Title: j.Title, Content: j.Content, Status: j.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Job{}, err
	}
	// if t.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
	// 	if err != nil {
	// 		return &Task{}, err
	// 	}
	// }
	return j, nil
}

func (j *Job) DeleteAJob(db *gorm.DB, jid uint64, aid uint32) (int64, error) {

	db = db.Debug().Model(&Job{}).Where("id = ? and author_id = ?", jid, aid).Take(&Job{}).Delete(&Job{})

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, echo.NewHTTPError(404, "Job not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
