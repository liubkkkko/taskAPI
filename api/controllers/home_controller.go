package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (server *Server) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Welcome To This Awesome API")
}

func (server *Server) TestRout(c echo.Context) error {
	return c.JSON(http.StatusOK, "Sak my balls 2 days")
}
