package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Author struct {
	ID         uint64       `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Nickname   string       `gorm:"column:nickname;size:255;not null;unique;" json:"nickname"`
	Email      string       `gorm:"column:email;size:100;not null;unique;" json:"email"`
	Password   string       `gorm:"column:password;size:100;not null;" json:"password"`
	Role       string       `gorm:"column:role;size:100;not null;default:'user'" json:"role"` //in the feature add interface to switch role
	Jobs       []Job        `gorm:"foreignKey:author_id"`
	Workspaces []*Workspace `gorm:"many2many:author_workspace;"`
	CreatedAt  time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (a *Author) BeforeSave(db *gorm.DB) error {
	hashedPassword, err := Hash(a.Password)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

func (a *Author) Prepare() {
	a.ID = 0
	a.Nickname = html.EscapeString(strings.TrimSpace(a.Nickname))
	a.Email = html.EscapeString(strings.TrimSpace(a.Email))
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
}

func (a *Author) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if a.Nickname == "" {
			return errors.New("required Nickname")
		}
		if a.Password == "" {
			return errors.New("required Password")
		}
		if a.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(a.Email); err != nil {
			return errors.New("invalid Email")
		}

		return nil
	case "login":
		if a.Password == "" {
			return errors.New("required Password")
		}
		if a.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(a.Email); err != nil {
			return errors.New("invalid Email")
		}
		return nil

	default:
		if a.Nickname == "" {
			return errors.New("required Nickname")
		}
		if a.Password == "" {
			return errors.New("required Password")
		}
		if a.Email == "" {
			return errors.New("required Email")
		}
		if err := checkmail.ValidateFormat(a.Email); err != nil {
			return errors.New("invalid Email")
		}
		return nil
	}
}

func (a *Author) SaveAuthors(db *gorm.DB) (*Author, error) {

	err := db.Debug().Create(&a).Error
	if err != nil {
		return &Author{}, err
	}
	return a, nil
}

func (a *Author) FindAllAuthors(db *gorm.DB) (*[]Author, error) {
	var authors []Author
	err := db.Debug().Preload("Workspaces").Preload("Jobs").Model(&Author{}).Limit(100).Find(&authors).Error
	if err != nil {
		return &[]Author{}, err
	}
	return &authors, err
}

func (a *Author) FindAuthorByIDForWorkspace(db *gorm.DB, aid uint32) error {
    err := db.Debug().Model(&Author{}).Where("id = ?", aid).First(a).Error
    if err != nil {
        return err
    }
    err = db.Debug().
        Model(&Author{}).
        Where("id = ?", a.ID).
        Preload("Workspaces.Jobs").
        Preload("Workspaces.Authors").
        Preload("Workspaces").
        Find(a).Error
    if err != nil {
        return err
    }
    return nil
}

func (a *Author) FindAuthorsByID(db *gorm.DB, uid uint32) (*Author, error) {
	err := db.Debug().Model(Author{}).Preload("Workspaces").Preload("Jobs").Where("id = ?", uid).Take(&a).Error
	if err != nil {
		return &Author{}, err
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return &Author{}, errors.New("user Not Found")
	}
	return a, err
}

func (a *Author) FindAuthorsByEmail(db *gorm.DB, email string) (*Author, error) {
	err := db.Debug().Model(Author{}).Where("email = ?", email).Take(&a).Error
	if err != nil {
		return &Author{}, err
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return &Author{}, errors.New("user Not Found")
	}
	return a, err
}

func (a *Author) UpdateAuthors(db *gorm.DB, uid uint32) (*Author, error) {
	// Хешування пароля перед збереженням
	err := a.BeforeSave(db)
	if err != nil {
		log.Fatal(err)
	}

	// Виконання оновлення
	db = db.Debug().Model(&Author{}).Where("id = ?", uid).Updates(map[string]interface{}{
		"password":   a.Password,
		"nickname":   a.Nickname,
		"email":      a.Email,
		"role":       a.Role,
		"updated_at": time.Now(),
	})
	if db.Error != nil {
		return &Author{}, db.Error
	}

	// Отримання оновленого запису
	err = db.Debug().Model(&Author{}).Where("id = ?", uid).Take(&a).Error
	if err != nil {
		return &Author{}, err
	}
	return a, nil
}


func (a *Author) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Author{}).Where("id = ?", uid).Take(&Author{}).Delete(&Author{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
