package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
	"github.com/liubkkkko/firstAPI/api/utils/formaterror"
)

// 🔹 LOGIN
func (server *Server) Login(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	var authorReq models.Author
	if err := json.Unmarshal(body, &authorReq); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	authorReq.Prepare()
	if err := authorReq.Validate("login"); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// ✅ Авторизація користувача
	author, accessToken, refreshToken, err := server.SignIn(authorReq.Email, authorReq.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		return c.JSON(http.StatusUnauthorized, formattedError)
	}

	// ✅ Встановлюємо access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(auth.AccessTokenTTL),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}
	c.SetCookie(accessCookie)

	// ✅ Встановлюємо refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(auth.RefreshTokenTTL),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}
	c.SetCookie(refreshCookie)

	// ✅ Очищаємо пароль
	author.Password = ""

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": author,
	})
}

// 🔹 REFRESH (rotation)
func (server *Server) Refresh(c echo.Context) error {
	refreshToken := auth.ExtractToken(c)
	if refreshToken == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
	}

	ctx := context.Background()
	userIdStr, err := tokenstorage.RedisClient.Get(ctx, refreshToken).Result()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
	}

	uid, err := strconv.Atoi(userIdStr)
	if err != nil || uid == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh mapping"})
	}

	newAccess, err := auth.CreateAccessToken(uint32(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "can't create access token"})
	}

	// створюємо новий jti і новий refresh (rotation)
	newJTI := uuid.New().String()
	newRefreshSigned, _, err := auth.CreateRefreshTokenWithJTI(uint32(uid), newJTI)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "can't create refresh token"})
	}
	// returnedJTI == newJTI (можна ігнорувати або перевіряти)

	// зберігаємо новий refresh у Redis
	if err := tokenstorage.RedisClient.Set(ctx, newRefreshSigned, uid, auth.RefreshTokenTTL).Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "can't save new refresh token"})
	}

	// видаляємо старий refresh токен
	_ = tokenstorage.RedisClient.Del(ctx, refreshToken).Err()

	// встановлюємо новий refresh cookie
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = newRefreshSigned
	cookie.Expires = time.Now().Add(auth.RefreshTokenTTL)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"access_token": newAccess})
}

// 🔹 LOGOUT (очищує лише refresh токен для цієї сесії)
func (server *Server) Logout(c echo.Context) error {
	refreshToken := auth.ExtractToken(c)
	ctx := context.Background()

	if refreshToken != "" {
		_ = tokenstorage.RedisClient.Del(ctx, refreshToken).Err()
	}

	expired := new(http.Cookie)
	expired.Name = "refresh_token"
	expired.Value = ""
	expired.Expires = time.Unix(0, 0)
	expired.Path = "/"
	expired.HttpOnly = true
	expired.Secure = true
	expired.SameSite = http.SameSiteNoneMode
	c.SetCookie(expired)

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// 🔹 SIGN-IN логіка (перевірка користувача, створення токенів)
func (server *Server) SignIn(email, password string) (models.Author, string, string, error) {
	var author models.Author
	err := server.DB.Debug().Model(models.Author{}).Where("email = ?", email).Take(&author).Error
	if err != nil {
		return author, "", "", err
	}

	if err := models.VerifyPassword(author.Password, password); err != nil {
		return author, "", "", err
	}

	ctx := context.Background()

	// === ЗАЛИШАЮ ТУТ ПРОСТУ ЛОГІКУ: дозволяємо multi-device (не блокуємо наявні refresh) ===
	// Якщо хочеш — можемо ввести обмеження на кількість сесій або логіку видалення старих сесій.

	// ✅ Створюємо access token
	accessToken, err := auth.CreateAccessToken(uint32(author.ID))
	if err != nil {
		return author, "", "", err
	}

	// ✅ створюємо refresh токен з унікальним jti
	jti := uuid.New().String()
	refreshSigned, returnedJTI, err := auth.CreateRefreshTokenWithJTI(uint32(author.ID), jti)
	if err != nil {
		return author, "", "", err
	}
	_ = returnedJTI // поки що не використовуємо окремо

	// ✅ Зберігаємо refresh токен у Redis (key = signedToken -> userID)
	if err := tokenstorage.RedisClient.Set(ctx, refreshSigned, author.ID, auth.RefreshTokenTTL).Err(); err != nil {
		return author, "", "", err
	}

	return author, accessToken, refreshSigned, nil
}

// 🔹 Отримати користувача за токеном
func (server *Server) IdIfYouHaveToken(c echo.Context) error {
	token := auth.ExtractToken(c)
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing token"})
	}

	ctx := context.Background()
	id, err := tokenstorage.RedisClient.Get(ctx, token).Result()
	if err != nil {
		if uid, err2 := auth.ExtractTokenID(c); err2 == nil && uid != 0 {
			author := models.Author{}
			userGotten, err := author.FindAuthorsByID(server.DB, uint32(uid))
			if err != nil {
				return c.JSON(http.StatusBadRequest, err)
			}
			return c.JSON(http.StatusOK, userGotten)
		}
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
