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
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().RequestURI, userCreated.ID))
	return c.JSON(http.StatusCreated, userCreated)
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
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(id)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	fmt.Println(body)
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	fmt.Println(user)
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	fmt.Println(tokenID)
	if tokenID != uint32(id) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	fmt.Println(user)
	updatedUser, err := user.UpdateAUser(server.DB, uint32(id))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	fmt.Println(updatedUser)
	return c.JSON(http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(c echo.Context) error {

	user := models.User{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	if tokenID != 0 && tokenID != uint32(id) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	_, err = user.DeleteAUser(server.DB, uint32(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", id))
	return c.JSON(http.StatusNoContent, "")
}
