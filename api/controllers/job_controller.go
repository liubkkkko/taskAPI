package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) CreateJob(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	job := models.Job{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	aid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	job.AuthorID = uint64(aid)
	job.Prepare()
	err = job.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	if aid != uint32(job.AuthorID) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	jobCreated, err := job.SaveJob(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().URL.Path, jobCreated.ID))
	return c.JSON(http.StatusCreated, jobCreated)
}

func (server *Server) GetJobs(c echo.Context) error {

	job := models.Job{}

	jobs, err := job.FindAllJob(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, jobs)
}

func (server *Server) GetJob(c echo.Context) error {

	tid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	job := models.Job{}

	taskReceived, err := job.FindJobByID(server.DB, uint64(tid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, taskReceived)
}


func (server *Server) GetJobsByWorkspaceId(c echo.Context) error {
	wid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	workspace := models.Workspace{}

	err = workspace.FindJobsByWorkspaceId(server.DB, uint32(wid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	jobs := workspace.Jobs
	return c.JSON(http.StatusOK, jobs)
}

func (server *Server) UpdateJob(c echo.Context) error {

	// Check if the job id is valid
	tId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the job exist
	job := models.Job{}
	err = server.DB.Debug().Model(models.Job{}).Where("id = ?", tId).Take(&job).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("post not found"))
	}

	// If a user attempt to update a job not belonging to him
	if uid != uint32(job.AuthorID) {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	// Read the data posted
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// Start processing the request data
	jobUpdate := models.Job{}
	err = json.Unmarshal(body, &jobUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != uint32(job.AuthorID) {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	jobUpdate.Prepare()
	err = jobUpdate.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	jobUpdate.ID = job.ID //this is important to tell the model the job id to update, the other update field are set above

	jobUpdated, err := jobUpdate.UpdateAJob(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	return c.JSON(http.StatusOK, jobUpdated)
}

func (server *Server) DeleteJob(c echo.Context) error {

	// Is a valid job id given to us?
	tid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the job exist
	job := models.Job{}
	err = server.DB.Debug().Model(models.Job{}).Where("id = ?", tid).Take(&job).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
	}

	// Is the authenticated user, the owner of this job?
	if uid != uint32(job.AuthorID) {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	_, err = job.DeleteAJob(server.DB, tid, uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", tid))
	return c.JSON(http.StatusNoContent, "")
}
