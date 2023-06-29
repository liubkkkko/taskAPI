package seed

import (
	"log"

	"github.com/liubkkkko/firstAPI/api/models"
	"gorm.io/gorm"
)

var users = []models.User{
	{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

var tasks = []models.Task{
	{
		Title:   "Task 1",
		Content: "Doing something interesting 1",
		Status:  "Created",
	},
	{
		Title:   "Task 2",
		Content: "Doing something interesting 2",
		Status:  "In proces",
	},
}

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
		&models.Post{},
		&models.User{},
		&models.Task{},
		&models.Job{},
		&models.Author{},
		&models.Workspace{},
	); err != nil {
		log.Fatal("failed to drop table")
	}

	if err = db.Debug().AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Task{},
		&models.Author{},
		&models.Workspace{},
		&models.Job{},
	); err != nil {
		log.Fatal("failed to drop table")
	}


	for i := range users {

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

		////////////////////////////////////////////////////////////////

		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		tasks[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Task{}).Create(&tasks[i]).Error
		if err != nil {
			log.Fatalf("cannot seed task table: %v", err)
		}
	}
}
