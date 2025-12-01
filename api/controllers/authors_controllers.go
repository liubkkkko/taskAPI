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
	"github.com/liubkkkko/firstAPI/api/responses"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) CreateAuthor(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
	}
	author := models.Author{}
	err = json.Unmarshal(body, &author)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	author.Prepare()
	err = author.Validate("")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	authorCreated, err := author.SaveAuthors(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().RequestURI, authorCreated.ID))
	return c.JSON(http.StatusCreated, authorCreated)
}

func (server *Server) GetAuthors(c echo.Context) error {
	author := models.Author{}
	authors, err := author.FindAllAuthors(server.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, authors)
}

func (server *Server) GetAuthor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("DEBUG: GetAuthor - parse id error:", err)
		return c.JSON(http.StatusNoContent, err)
	}
	author := models.Author{}
	userGotten, err := author.FindAuthorsByID(server.DB, uint32(id))
	if err != nil {
		fmt.Println("DEBUG: GetAuthor - FindAuthorsByID error:", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, userGotten)
}

func (server *Server) UpdateAuthor(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	author := models.Author{}
	err = json.Unmarshal(body, &author)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	tokenID, err := auth.ExtractTokenID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
	}
	if tokenID != uint32(id) {
		return c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	author.Prepare()
	err = author.Validate("update")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	updatedAuthor, err := author.UpdateAuthors(server.DB, uint32(id))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusInternalServerError, formattedError)
	}
	return c.JSON(http.StatusOK, updatedAuthor)
}

func (server *Server) DeleteAuthor(c echo.Context) error {

	author := models.Author{}
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
	_, err = author.DeleteAUser(server.DB, uint32(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("Entity", fmt.Sprintf("%d", id))
	return c.JSON(http.StatusNoContent, "")
}
