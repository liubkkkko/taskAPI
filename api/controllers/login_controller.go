package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return echo.NewHTTPError(http.StatusUnprocessableEntity, "Please provide valid credentials")

	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, err)
		// return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, formattedError)
		// responses.ERROR(c.Response(), http.StatusUnprocessableEntity, formattedError)
		// return
	}
	// responses.JSON(c.Response(), http.StatusOK, token)
	return c.JSON(http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
