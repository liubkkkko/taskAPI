package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
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
	// w.Author = Author{}
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
	for i, author := range w.Authors {
		authorIDs[i] = author.ID
	}
	return authorIDs

}

func (w *Workspace) CheckIfYouAuthor(aid uint64) error {
	allUsersId := w.GetAllAuthorsId()
	authorised := false
	for i := range allUsersId {
		if aid == allUsersId[i] {
			authorised = true
		}
	}
	if !authorised {
		return echo.ErrUnauthorized
	}
	return nil

}

func (w *Workspace) SaveWorkspace(db *gorm.DB) (*Workspace, error) {
	err := db.Debug().Model(&Workspace{}).Create(&w).Error
	if err != nil {
		return &Workspace{}, err
	}
	//IN THE FEAUTHURE ADD CHECK TO WRITE SAME DATA
	// if w.ID != 0 {
	// 	err = db.Debug().Model(&Task{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
	// 	if err != nil {
	// 		return &Task{}, err
	// 	}
	// }
	return w, nil
}

func (w *Workspace) FindAllWorkspaces(db *gorm.DB) (*[]Workspace, error) {
	workspace := []Workspace{}
	err := db.Debug().Model(&Workspace{}).Limit(100).Find(&workspace).Error
	if err != nil {
		return &[]Workspace{}, err
	}
	// if len(workspace) > 0 {
	// 	for i := range workspace {
	// 		err := db.Debug().Model(&User{}).Where("id = ?", tasks[i].AuthorID).Take(&tasks[i].Author).Error
	// 		if err != nil {
	// 			return &[]Task{}, err
	// 		}
	// 	}
	// }
	return &workspace, nil
}

func (w *Workspace) AddAuthorsToWorkspace(db *gorm.DB, wid uint32) error {
	err := db.Debug().Model(&Author{}).Where("id = ?", wid).Take(&w.Authors).Error
	if err != nil {
		return err
	}
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
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, echo.NewHTTPError(404, "Task not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
