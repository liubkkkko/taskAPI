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
		Title:   "Task1",
		Content: "heh work",
	},
	{
		Title:   "Task2",
		Content: "heh workkk",
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
	if err = db.Debug().Migrator().DropTable(
		&models.Job{},
		&models.Author{},
		&models.Workspace{},
	); err != nil {
		log.Fatal("failed to drop table")
	}

	if err = db.Debug().AutoMigrate(
		&models.Author{},
		&models.Workspace{},
		&models.Job{},
	); err != nil {
		log.Fatal("failed to drop table")
	}

	for i := range authors {

		err := db.Debug().Model(&models.Workspace{}).Create(&workspaces[i]).Error
		if err != nil {
			log.Fatalf("cannot seed workspace table: %v", err)
		}

		err = db.Debug().Model(&models.Author{}).Create(&authors[i]).Error
		if err != nil {
			log.Fatalf("cannot seed authors table: %v", err)
		}

		err = db.Debug().Model(&models.Job{}).Create(&jobs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed jobs table: %v", err)
		}

		jobs[i].WorkspaceID = workspaces[i].ID
		jobs[i].AuthorID = authors[i].ID
	}
}
