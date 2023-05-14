package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) CreateTask(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	task := models.Task{}
	err = json.Unmarshal(body, &task)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	task.Prepare()
	err = task.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	} 
	if uid != task.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	postCreated, err := task.SaveTask(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Lacation", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().URL.Path, postCreated.ID))
	return c.JSON(http.StatusCreated, postCreated)
}

func (server *Server) GetTasks(c echo.Context) error {

	task := models.Task{}

	posts, err := task.FindAllTasks(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, posts)
}
