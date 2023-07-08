package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

func (server *Server) Login(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	author := models.Author{}
	err = json.Unmarshal(body, &author)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	author.Prepare()
	err = author.Validate("login")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	token, err := server.SignIn(author.Email, author.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, formattedError)
	}
	return c.JSON(http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {
	author := models.Author{}

	err := server.DB.Debug().Model(models.Author{}).Where("email = ?", email).Take(&author).Error
	if err != nil {
		return "", err
	}
	//for test
	fmt.Println("author password", author.Password, "password", password)
	err = models.VerifyPassword(author.Password, password)
	if err != nil {
		return "bad login data", err
	}
	return auth.CreateToken(uint32(author.ID))
}

//&& err == bcrypt.ErrMismatchedHashAndPassword
