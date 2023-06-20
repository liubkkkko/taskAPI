package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Workspace struct {
	ID          uint64    `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"column:name;size:255;not null;" json:"name"`
	Description string    `gorm:"column:description;size:255;not null;" json:"description"`
	Status      string    `gorm:"column:status;size:100;not null;default:'created'" json:"status"`
	Jobs        []Job     `gorm:"foreignKey:workspace_id"`
	Authors     []*Author `gorm:"many2many:author_workspace;"`
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

func (w *Workspace) CheckIfYouAuthor(aid uint64) error {
	for _, id := range w.GetAllAuthorsId() {
		if id == aid {
			return nil
		}
	}
	return echo.ErrUnauthorized

}

func (w *Workspace) SaveWorkspace(db *gorm.DB) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Create(&w).Error
	if err != nil {
		return &Workspace{}, err
	}
	fmt.Println("After save workspace", w.Authors[0], w)
	return w, nil
}

func (w *Workspace) FindAllWorkspaces(db *gorm.DB) (*[]Workspace, error) {
	workspace := []Workspace{}
	err := db.Debug().Model(&Workspace{}).Limit(100).Find(&workspace).Error
	if err != nil {
		return &[]Workspace{}, err
	}
	if len(workspace) > 0 {
		for i := range workspace {
			authorsIds := workspace[i].GetAllAuthorsId()
			for j := range authorsIds {
				author := &Author{}
				err := db.Debug().Model(&Author{}).Where("id = ?", authorsIds[j]).First(author).Error
				if err != nil {
					return &[]Workspace{}, err
				}
				workspace[i].Authors = append(workspace[i].Authors, author)
			}
		}
	}
	return &workspace, nil

}

func (w *Workspace) AddAuthorsToWorkspace(db *gorm.DB, aid uint32) error {
	err := db.Debug().Model(&Author{}).Where("id = ?", aid).Take(&w.Authors).Error
	if err != nil {
		return err
	}
	// err = db.Debug().Save(&w).Error
	// if err != nil {
	// 	log.Fatalf("error to save: %v", err)
	// }
	return nil
}

func (w *Workspace) FindWorkspaceByID(db *gorm.DB, pid uint64) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Where("id = ?", pid).Take(&w).Error
	if err != nil {
		return &Workspace{}, err
	}
	// if t.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
	// 	if err != nil {
	// 		return &Task{}, err
	// 	}
	// }
	return w, nil
}

// func (w *Workspace) FindWorkspaceByAuthorID(db *gorm.DB, aid uint64) (*Workspace, error) {
// 	err :=  db.Debug().Model(&Workspace{}).Where()
// 	err = db.Debug().Model(&Workspace{}).Where("id = ?", aid).Take(&w).Error
// 	if err != nil {
// 		return &Workspace{}, err
// 	}
// 	// if t.ID != 0 {
// 	// 	err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
// 	// 	if err != nil {
// 	// 		return &Task{}, err
// 	// 	}
// 	// }
// 	return w, nil
// }


func (w *Workspace) UpdateWorkspace(db *gorm.DB) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Where("id = ?", w.ID).Updates(Workspace{Name: w.Name, Description: w.Description, Status: w.Status, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Workspace{}, err
	}
	// if t.ID != 0 {
	// 	err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
	// 	if err != nil {
	// 		return &Task{}, err
	// 	}
	// }
	return w, nil
}

func (w *Workspace) DeleteAWorkspace(db *gorm.DB, tid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Workspace{}).Where("id = ? and author_id = ?", tid, uid).Take(&Workspace{}).Delete(&Workspace{})
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, echo.NewHTTPError(404, "Task not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
