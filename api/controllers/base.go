package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/liubkkkko/firstAPI/api/models"
)

type Server struct {
	DB *gorm.DB
	// Router *mux.Router
	Router *echo.Echo
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	// server.Router = mux.NewRouter()
	server.Router = echo.New()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	
	fmt.Println("Listening to port 8080")
	if err := server.Router.Start(addr); err != http.ErrServerClosed{
		log.Fatal(err)
	}
	// if err := server.Router.Start(addr); err!=http.ErrServerClosed{
	// 	log.Fatal(err)
	// }
	// log.Fatal(http.ListenAndServe(addr, server.Router))
}
