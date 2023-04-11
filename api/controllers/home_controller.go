package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (server *Server) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Welcome To This Awesome API")
}
