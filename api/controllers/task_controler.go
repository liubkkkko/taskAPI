package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

	tasks, err := task.FindAllTasks(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, tasks)
}

func (server *Server) GetTask(c echo.Context) error {

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	task := models.Task{}

	taskReceived, err := task.FindTaskByID(server.DB, uint64(tid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, taskReceived)
}

func (server *Server) UpdateTask(c echo.Context) error {

	// Check if the task id is valid
	tId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", tId).Take(&task).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("post not found"))
	}

	// If a user attempt to update a post not belonging to him
	if uid != task.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// Start processing the request data
	taskUpdate := models.Task{}
	err = json.Unmarshal(body, &taskUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != taskUpdate.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	taskUpdate.Prepare()
	err = taskUpdate.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	taskUpdate.ID = task.ID //this is important to tell the model the post id to update, the other update field are set above

	taskUpdated, err := taskUpdate.UpdateATask(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	return c.JSON(http.StatusOK, taskUpdated)
}

func (server *Server) DeleteTask(c echo.Context) error {

	// Is a valid task id given to us?
	tid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	

	// Check if the task exist
	task := models.Task{}
	err = server.DB.Debug().Model(models.Task{}).Where("id = ?", tid).Take(&task).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
	}

	// Is the authenticated user, the owner of this task?
	if uid != task.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	_, err = task.DeleteATask(server.DB, tid, uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", tid))
	return c.JSON(http.StatusNoContent, "")
}
