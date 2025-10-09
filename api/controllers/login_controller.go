package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Token  string        `json:"token"`
	Author models.Author `json:"author"`
}

// 🔹 LOGIN
func (server *Server) Login(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	author := models.Author{}
	if err := json.Unmarshal(body, &author); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	author.Prepare()
	if err := author.Validate("login"); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	loginResponse, err := server.SignIn(author.Email, author.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusUnauthorized, formattedError)
	}

	// ✅ створюємо cookie через Echo API
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = loginResponse.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode // для різних доменів

	c.SetCookie(cookie) // ✅ додає cookie у відповідь

	// повертаємо дані користувача (без токена)
	return c.JSON(http.StatusOK, loginResponse.Author)
}

// 🔹 LOGOUT
func (server *Server) Logout(c echo.Context) error {
	token := auth.ExtractToken(c)
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing token"})
	}

	ctx := context.Background()
	if err := tokenstorage.RedisClient.Del(ctx, token).Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Видаляємо cookie
	expired := new(http.Cookie)
	expired.Name = "access_token"
	expired.Value = ""
	expired.Expires = time.Unix(0, 0)
	expired.Path = "/"
	expired.HttpOnly = true
	expired.Secure = true
	c.SetCookie(expired)

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// 🔹 SIGN-IN логіка (перевірка користувача, запис токена)
func (server *Server) SignIn(email, password string) (*LoginResponse, error) {
	author := models.Author{}
	if err := server.DB.Debug().Model(models.Author{}).Where("email = ?", email).Take(&author).Error; err != nil {
		return nil, err
	}

	if err := models.VerifyPassword(author.Password, password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	token, err := auth.CreateToken(uint32(author.ID))
	if err != nil {
		return nil, errors.New("can't create token")
	}

	ctx := context.Background()
	if err := tokenstorage.RedisClient.Set(ctx, token, author.ID, time.Hour*24).Err(); err != nil {
		return nil, err
	}

	author.Password = "" // очищаємо перед відправкою

	return &LoginResponse{
		Token:  token,
		Author: author,
	}, nil
}

// 🔹 Отримати користувача, якщо є токен
func (server *Server) IdIfYouHaveToken(c echo.Context) error {
	token := auth.ExtractToken(c)
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing token"})
	}

	ctx := context.Background()
	id, err := tokenstorage.RedisClient.Get(ctx, token).Result()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
	}

	Id, err := strconv.Atoi(id)
	if err != nil || Id == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	author := models.Author{}
	userGotten, err := author.FindAuthorsByID(server.DB, uint32(Id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, userGotten)
}