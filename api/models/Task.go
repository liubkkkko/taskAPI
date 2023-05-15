package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Task struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Status    string    `gorm:"size:255;not null;" json:"status"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (t *Task) Prepare() {
	t.ID = 0
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
	t.Content = html.EscapeString(strings.TrimSpace(t.Content))
	t.Status = html.EscapeString(strings.TrimSpace(t.Status))
	t.Author = User{}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}

func (t *Task) Validate() error {

	if t.Title == "" {
		return echo.NewHTTPError(418, "required Title")
	}
	if t.Content == "" {
		return echo.NewHTTPError(418, "required Content")
	}
	if t.AuthorID < 1 {
		return echo.NewHTTPError(418, "required Author")
	}
	if t.Status == "" {
		return echo.NewHTTPError(418, "required Status")
	}
	return nil
}

func (t *Task) SaveTask(db *gorm.DB) (*Task, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&t).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&Task{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) FindAllTasks(db *gorm.DB) (*[]Task, error) {
	var err error
	tasks := []Task{}
	err = db.Debug().Model(&Post{}).Limit(100).Find(&tasks).Error
	if err != nil {
		return &[]Task{}, err
	}
	if len(tasks) > 0 {
		for i := range tasks {
			err := db.Debug().Model(&User{}).Where("id = ?", tasks[i].AuthorID).Take(&tasks[i].Author).Error
			if err != nil {
				return &[]Task{}, err
			}
		}
	}
	return &tasks, nil
}

func (t *Task) FindTaskByID(db *gorm.DB, pid uint64) (*Task, error) {
	var err error
	err = db.Debug().Model(&Task{}).Where("id = ?", pid).Take(&t).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) UpdateATask(db *gorm.DB) (*Task, error) {

	var err error

	err = db.Debug().Model(&Task{}).Where("id = ?", t.ID).Updates(Task{Title: t.Title, Content: t.Content, Status: t.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Task{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &Task{}, err
		}
	}
	return t, nil
}

func (t *Task) DeleteATask(db *gorm.DB, tid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Task{}).Where("id = ? and author_id = ?", tid, uid).Take(&Task{}).Delete(&Task{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, echo.NewHTTPError(404, "Task not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
