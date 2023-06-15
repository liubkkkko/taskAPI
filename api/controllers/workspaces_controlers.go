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

func (server *Server) CreateWorspace(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
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
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Lacation", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().URL.Path, postCreated.ID))
	return c.JSON(http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(c echo.Context) error {

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, posts)
}

func (server *Server) GetPost(c echo.Context) error {

	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, uint64(pid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(c echo.Context) error {

	// Check if the post id is valid
	pid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("post not found"))
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	return c.JSON(http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(c echo.Context) error {

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", pid))
	return c.JSON(http.StatusNoContent, "")
}
