package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (server *Server) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Welcome To This Awesome API")
}

func (server *Server) TestRout(c echo.Context) error {
	return c.JSON(http.StatusOK, "created new routs")
}

func (server *Server) Health(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "status": "ok",
    })
}
