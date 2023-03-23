package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/responses"
)

func (server *Server) Home(c echo.Context) {
	responses.JSON(c.Response(), http.StatusOK, "Welcome To This Awesome API")

}