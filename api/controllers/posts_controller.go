package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) CreatePost(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("unauthorized"))
		// return
	}
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, formattedError)
		// return
	}
	c.Response().Header().Set("Lacation", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().URL.Path, postCreated.ID))
	// w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	// responses.JSON(c.Response(), http.StatusCreated, postCreated)
	return c.JSON(http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(c echo.Context) error {

	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, err)
		// return
	}
	return c.JSON(http.StatusOK, posts)
	// responses.JSON(c.Response(), http.StatusOK, posts)
}

func (server *Server) GetPost(c echo.Context) error {

	vars := mux.Vars(c.Request())
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}
	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, err)
		// return
	}
	// responses.JSON(c.Response(), http.StatusOK, postReceived)
	return c.JSON(http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(c echo.Context) error {

	vars := mux.Vars(c.Request())

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("unauthorized"))
		// return
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("post not found"))
		// responses.ERROR(c.Response(), http.StatusNotFound, errors.New("post not found"))
		// return
	}

	// If a user attempt to update a post not belonging to him
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("unauthorized"))
		// return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}

	// Start processing the request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != postUpdate.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("Unauthorized"))
		// return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}

	postUpdate.ID = post.ID //this is important to tell the model the post id to update, the other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, formattedError)
		// return
	}
	// responses.JSON(c.Response(), http.StatusOK, postUpdated)
	return c.JSON(http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(c echo.Context) error {

	vars := mux.Vars(c.Request())

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("unauthorized"))
		// return
	}

	// Check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
		// 	responses.ERROR(c.Response(), http.StatusNotFound, errors.New("Unauthorized"))
		// 	return
	}

	// Is the authenticated user, the owner of this post?
	if uid != post.AuthorID {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("Unauthorized"))
		// return
	}
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", pid))
	// w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	// responses.JSON(c.Response(), http.StatusNoContent, "")
	return c.JSON(http.StatusNoContent, "")
}
