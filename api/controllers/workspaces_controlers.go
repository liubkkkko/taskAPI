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

func (server *Server) CreateWorkspace(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	workspace := models.Workspace{}
	err = json.Unmarshal(body, &workspace)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	workspace.Prepare()
	err = workspace.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	aid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	_, err = workspace.AddAuthorToWorkspace(server.DB, aid, uint32(workspace.ID))
	if err != nil {
		return c.JSON(http.StatusFailedDependency, err)
	}
	err = workspace.CheckIfYouAuthor(server.DB, uint64(aid))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	workspaceCreated, err := workspace.SaveWorkspace(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().URL.Path, workspaceCreated.ID))
	return c.JSON(http.StatusCreated, workspaceCreated)
}

func (server *Server) GetWorkspacesByAuthorId(c echo.Context) error {
	aid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	author := models.Author{}

	err = author.FindAuthorByIDForWorkspace(server.DB, uint32(aid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	workspaces := author.Workspaces
	return c.JSON(http.StatusOK, workspaces)
}

func (server *Server) CheckIfIAuthor(c echo.Context) error {
	wid, err := strconv.Atoi(c.FormValue("wid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	aid, err := strconv.Atoi(c.FormValue("aid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Check if the workspace exist
	workspace := models.Workspace{}
	err = server.DB.Debug().Model(models.Workspace{}).Where("id = ?", wid).Take(&workspace).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("workspace not found"))
	}

	err = workspace.CheckIfYouAuthor(server.DB, uint64(aid))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("check is not valid"))
	}

	return nil
}

func (server *Server) GetWorkspaces(c echo.Context) error {

	workspace := models.Workspace{}

	workspaces, err := workspace.FindAllWorkspaces(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, workspaces)

}

func (server *Server) GetWorkspace(c echo.Context) error {
	wid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	workspace := models.Workspace{}

	workspaceReceived, err := workspace.FindWorkspaceByID(server.DB, uint64(wid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, workspaceReceived)
}

func (server *Server) AddOneMoreAuthorToWorkspace(c echo.Context) error {

	wid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	aid, err := strconv.Atoi(c.FormValue("aid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	workspace := models.Workspace{}

	workspacePlusOneAuthor, err := workspace.AddAuthorToWorkspace(server.DB, uint32(aid), uint32(wid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, workspacePlusOneAuthor)
}

func (server *Server) UpdateWorkspace(c echo.Context) error {

	// Check if the workspace id is valid
	wid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//CHeck if the auth token is valid and  get the user id from it
	aid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the workspace exist
	workspace := models.Workspace{}
	err = server.DB.Debug().Model(models.Workspace{}).Where("id = ?", wid).Take(&workspace).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("post not found"))
	}

	// If a user attempt to update a workspace not belonging to him
	err = workspace.CheckIfYouAuthor(server.DB, uint64(aid))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Read the data posted
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// Start processing the request data
	workspaceUpdate := models.Workspace{}
	err = json.Unmarshal(body, &workspaceUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	workspaceUpdate.Prepare()
	err = workspaceUpdate.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	workspaceUpdate.ID = workspace.ID //this is important to tell the model the workspace id to update, the other update field are set above

	workspaceUpdated, err := workspaceUpdate.UpdateWorkspace(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	return c.JSON(http.StatusOK, workspaceUpdated)
}

func (server *Server) DeleteWorkspace(c echo.Context) error {

	// Is a valid post id given to us?
	wid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Is this user authenticated?
	aid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the workspace exist
	workspace := models.Workspace{}
	err = server.DB.Debug().Model(models.Workspace{}).Where("id = ?", wid).Take(&workspace).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
	}

	// Is the authenticated user, the owner of this post?
	err = workspace.CheckIfYouAuthor(server.DB, uint64(aid))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized2"))
	}

	_, err = workspace.DeleteAWorkspace(server.DB, wid, aid)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", wid))
	return c.JSON(http.StatusNoContent, "")
}
