package controllers

import (
    "fmt"
    "log"
    "net/http"

    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "go.uber.org/zap"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Server struct {
    DB     *gorm.DB
    Router *echo.Echo
    Logger *zap.Logger
}


func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
    var err error
    dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
    server.DB, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
    if err != nil {
        fmt.Printf("Cannot connect to %s database", DbDriver)
        log.Fatal("This is the error:", err)
    } else {
        fmt.Printf("We are connected to the %s database", DbDriver)
    }

    server.Router = echo.New()

    server.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
        AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
        AllowCredentials: true,
    }))

logger, err := zap.NewProduction()
if err != nil {
    log.Fatal(err)
}
server.Logger = logger

server.Router.Use(middleware.RequestLoggerWithConfig(
    middleware.RequestLoggerConfig{
        LogURI:    true,
        LogStatus: true,
        LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
            server.Logger.Info("request",
                zap.String("uri", v.URI),
                zap.Int("status", v.Status),
            )
            return nil
        },
    },
))


    server.initializeRoutes()
}

func (server *Server) Run(addr string) {
    fmt.Println("Listening (HTTP) on", addr)
    if err := server.Router.Start(addr); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}