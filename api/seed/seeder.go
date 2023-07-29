package seed

import (
	"log"

	"github.com/liubkkkko/firstAPI/api/models"
	"gorm.io/gorm"
)

var authors = []models.Author{
	{
		Nickname: "jinzhu",
		Email:    "jinzhu@gmail.com",
		Password: "jinzhu123",
	},
	{
		Nickname: "liubkkkk0",
		Email:    "liubkkkk0@gmail.com",
		Password: "liubkkkk0322",
	},
}

var jobs = []models.Job{
	{
		Title:       "Task1",
		Content:     "heh work",
		AuthorID:    1,
		WorkspaceID: 1,
	},
	{
		Title:       "Task2",
		Content:     "heh workkk",
		AuthorID:    2,
		WorkspaceID: 2,
	},
}

var workspaces = []models.Workspace{
	{
		Name:        "Workspace1",
		Description: "heh work",
	},

	{
		Name:        "Workspace2",
		Description: "heh workkkk",
	},
}

func Load(db *gorm.DB) {
	var err error

	// Міграція таблиць
	if err = db.Debug().AutoMigrate(
		&models.Author{},
		&models.Workspace{},
		&models.Job{},
	); err != nil {
		log.Fatal("failed to migrate tables")
	}

	// Перевірка, чи існують дані в таблицях
	var authorsCount int64
	if err = db.Model(&models.Author{}).Count(&authorsCount).Error; err != nil {
		log.Fatalf("failed to check authors table: %v", err)
	}

	var workspacesCount int64
	if err = db.Model(&models.Workspace{}).Count(&workspacesCount).Error; err != nil {
		log.Fatalf("failed to check workspaces table: %v", err)
	}

	var jobsCount int64
	if err = db.Model(&models.Job{}).Count(&jobsCount).Error; err != nil {
		log.Fatalf("failed to check jobs table: %v", err)
	}

	// If we don't have data - start to seed
	if authorsCount == 0 && workspacesCount == 0 && jobsCount == 0 {
		for i := range authors {
			err = db.Debug().Model(&models.Workspace{}).Create(&workspaces[i]).Error
			if err != nil {
				log.Fatalf("cannot seed workspace table: %v", err)
			}

			err = db.Debug().Model(&models.Author{}).Create(&authors[i]).Error
			if err != nil {
				log.Fatalf("cannot seed authors table: %v", err)
			}

			// append author to workspace
			if err := db.Debug().Model(&workspaces[i]).Association("Authors").Append(&authors[i]); err != nil {
				log.Fatalf("cannot append author to workspace: %v", err)
			}

			err = db.Debug().Model(&models.Job{}).Create(&jobs[i]).Error
			if err != nil {
				log.Fatalf("cannot seed jobs table: %v", err)
			}
		}
	} else {
		log.Println("Skipping data seeding, tables already have data.")
	}
}
