package controllers

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *echo.Echo
	// redisClient *redis.Client
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	//connect to postgres
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	//create new instance router
	server.Router = echo.New()

	// //connect to redis
	// server.redisClient = redis.NewClient(&redis.Options{
	// 	Addr:     RedisAddr,
	// 	Password: RedisPassword,
	// 	DB:       RedisDb,
	// })

	// // Check connect to Redis
	// pong, err := server.redisClient.Ping(context.Background()).Result()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Redis: %s", err)
	// }
	// fmt.Printf("Connected to Redis: %s\n", pong)

	//initialize logger
	logger, _ := zap.NewProduction()
	server.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	if err := server.Router.Start(addr); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
