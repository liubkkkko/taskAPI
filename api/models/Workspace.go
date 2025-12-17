package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

type Workspace struct {
	ID          uint64    `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"column:name;size:255;not null;unique;" json:"name"`
	Description string    `gorm:"column:description;size:255;not null;" json:"description"`
	Status      string    `gorm:"column:status;size:100;not null;default:'created'" json:"status"`
	Jobs        []Job     `gorm:"foreignKey:workspace_id;constraint:OnDelete:CASCADE;"`
	Authors     []*Author `gorm:"many2many:author_workspace;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (w *Workspace) Prepare() {
	w.ID = 0
	w.Name = html.EscapeString(strings.TrimSpace(w.Name))
	w.Description = html.EscapeString(strings.TrimSpace(w.Description))
	w.Status = html.EscapeString(strings.TrimSpace(w.Status))
	w.Jobs = []Job{}
	w.Authors = []*Author{}
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
}

func (w *Workspace) Validate() error {
	if w.Name == "" {
		return echo.NewHTTPError(418, "required Name")
	}
	if w.Description == "" {
		return echo.NewHTTPError(418, "required Description")
	}
	if w.Status == "" {
		return echo.NewHTTPError(418, "required Status")
	}
	return nil
}

func (w *Workspace) GetAllAuthorsId() []uint64 {
	authorIDs := make([]uint64, len(w.Authors))
	for _, author := range w.Authors {
		authorIDs = append(authorIDs, author.ID)
	}
	return authorIDs
}

func (w *Workspace) CheckIfYouAuthor(db *gorm.DB, aid uint64) error {

	if w.ID == 0 {
		authorIDs := make([]uint64, len(w.Authors)) // create array IDs authors
		for i, author := range w.Authors {
			authorIDs[i] = author.ID
		}
		for _, id := range authorIDs {
			if id == aid {
				return nil
			}
		}
		return echo.ErrUnauthorized
	} else {
		var workspace Workspace
		if err := db.Preload("Authors").First(&workspace, w.ID).Error; err != nil {
			return err
		}
		authorIDs := make([]uint64, len(workspace.Authors))
		for i, author := range workspace.Authors {
			authorIDs[i] = author.ID
		}
		for _, id := range authorIDs {
			if id == aid {
				return nil
			}
		}
		return echo.ErrUnauthorized
	}
}

func (w *Workspace) FindAllWorkspaces(db *gorm.DB) (*[]Workspace, error) {
	var workspaces []Workspace
	err := db.Debug().Preload("Authors").Preload("Jobs").Limit(100).Find(&workspaces).Error
	if err != nil {
		return &[]Workspace{}, err
	}
	return &workspaces, nil
}

func (w *Workspace) SaveWorkspace(db *gorm.DB) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Create(&w).Error
	if err != nil {
		return &Workspace{}, err
	}
	return w, nil
}

func (w *Workspace) AddAuthorToWorkspace(db *gorm.DB, aid uint32, wid uint32) (*Workspace, error) {
	// find author by id
	var existingAuthor Author
	if err := db.First(&existingAuthor, aid).Error; err != nil {
		return nil, err
	}
	var existingWorkspace Workspace
	if wid != 0 {
		// find workspace by id
		if err := db.Debug().First(&existingWorkspace, wid).Error; err != nil {
			return nil, err
		}
		// append author to workspace
		if err := db.Debug().Model(&existingWorkspace).Association("Authors").Append(&existingAuthor); err != nil {
			return nil, err
		}
	} else { //for create new workspace
		err := db.Debug().Model(&Author{}).Where("id = ?", aid).Take(&w.Authors).Error
		if err != nil {
			return nil, err
		}
	}
	return &existingWorkspace, nil
}

func (w *Workspace) FindWorkspaceByID(db *gorm.DB, wid uint64) (*Workspace, error) {
	fmt.Println("in FindWorkspaceByID")
	err := db.Debug().Model(&Workspace{}).Preload("Authors").Where("id = ?", wid).Take(&w).Error
	if err != nil {
		return &Workspace{}, err
	}
	return w, nil
}

func (w *Workspace) UpdateWorkspace(db *gorm.DB) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Where("id = ?", w.ID).Updates(Workspace{Name: w.Name, Description: w.Description, Status: w.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Workspace{}, err
	}
	return w, nil
}

func (w *Workspace) DeleteAWorkspace(db *gorm.DB, wid uint64, aid uint32) (int64, error) {
	db = db.Debug().Model(&Workspace{}).Where("id = ?", wid).Take(&Workspace{}).Delete(&Workspace{})
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, echo.NewHTTPError(404, "Task not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (w *Workspace) FindJobsByWorkspaceId(db *gorm.DB, wid uint32) error {
	err := db.Debug().Model(&Workspace{}).Where("id = ?", wid).Take(&w).Error
	if err != nil {
		return err
	}
err = db.Debug().
	Model(&Workspace{}).
	Where("id = ?", w.ID).
	Preload("Jobs").
	Preload("Authors").
	First(w).Error
	if err != nil {
		return err
	}
	return nil
}
