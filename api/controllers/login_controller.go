package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
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

func (server *Server) Logout(c echo.Context) error {
	token := auth.ExtractToken(c)
	fmt.Println(token)
	ctx := context.Background()
	// delete token in Redis
	err := tokenstorage.RedisClient.Del(ctx, token).Err()
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) SignIn(email, password string) (string, error) {
	author := models.Author{}

	err := server.DB.Debug().Model(models.Author{}).Where("email = ?", email).Take(&author).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(author.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "bad login data", err
	}

	token, err := auth.CreateToken(uint32(author.ID))
	if err != nil {
		return "can't create token", err
	}
	ctx := context.Background()
	// save token in Redis temporary (24h)
	err = tokenstorage.RedisClient.Set(ctx, token, author.ID, time.Hour*24).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}
