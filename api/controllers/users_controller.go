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
	"github.com/liubkkkko/firstAPI/api/responses"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) CreateUser(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		return c.JSON(http.StatusInternalServerError, formattedError)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, formattedError)
		// return
	}
	c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().RequestURI, userCreated.ID))
	// w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	return c.JSON(http.StatusCreated, userCreated)
	// responses.JSON(c.Response(), http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(c echo.Context) error {
	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (server *Server) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNoContent, err)
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(c echo.Context) error {
	// vars := mux.Vars(c.Request())
	// uid, err := strconv.ParseUint(vars["id"], 10, 32)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("Unauthorized"))
		// return
	}
	if tokenID != uint32(id) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(id))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, formattedError)
		// return
	}
	// responses.JSON(c.Response(), http.StatusOK, updatedUser)
	return c.JSON(http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(c echo.Context) error {

	vars := mux.Vars(c.Request())

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
		// responses.ERROR(c.Response(), http.StatusBadRequest, err)
		// return
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("Unauthorized"))
		// return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		// return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
		// responses.ERROR(c.Response(), http.StatusInternalServerError, err)
		// return
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", uid))
	// c.Header().Set("Entity", fmt.Sprintf("%d", uid))
	// responses.JSON(c.Response(), http.StatusNoContent, "")
	return c.JSON(http.StatusNoContent, "")
}
