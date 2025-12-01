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
    author, accessToken, refreshToken, jti, err := server.SignIn(authorReq.Email, authorReq.Password)
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
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(accessCookie)

    // ✅ Встановлюємо refresh token cookie
    refreshCookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Expires:  time.Now().Add(auth.RefreshTokenTTL),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(refreshCookie)

    // ✅ Зберігаємо метадані сесії в Redis (ключ = refresh:{jti})
    meta := tokenstorage.SessionMeta{
        UserID:    int(author.ID),
        IP:        c.RealIP(),
        UserAgent: c.Request().Header.Get("User-Agent"),
        CreatedAt: time.Now().Unix(),
        Signed:    refreshToken, // опціонально для дебагу; можна не зберігати в продакшні
    }
    _ = tokenstorage.SaveSession(tokenstorage.RedisClient, jti, meta, auth.RefreshTokenTTL)

    // ✅ Очищаємо пароль
    author.Password = ""

    return c.JSON(http.StatusOK, map[string]interface{}{
        "user": author,
    })
}

// 🔹 REFRESH (rotation)
// ...existing code...
// 🔹 REFRESH (rotation)
func (server *Server) Refresh(c echo.Context) error {
    refreshToken := auth.ExtractToken(c)
    if refreshToken == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
    }

    // спочатку пробуємо витягти jti із токена
    oldJTI, err := auth.ExtractJTIFromString(refreshToken)
    ctx := context.Background()
    var uid int
    var meta tokenstorage.SessionMeta
    if err == nil && oldJTI != "" {
        // новий підхід: jti -> session
        sMeta, err2 := tokenstorage.GetSession(tokenstorage.RedisClient, oldJTI)
        if err2 == nil {
            meta = sMeta
            uid = sMeta.UserID
        }
    }

    // backward-compat: якщо не знайшли по jti, пробуємо стару схему (ключ = signed token)
    if uid == 0 {
        userIdStr, err := tokenstorage.RedisClient.Get(ctx, refreshToken).Result()
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
        }
        uid64, err := strconv.Atoi(userIdStr)
        if err != nil || uid64 == 0 {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh mapping"})
        }
        uid = uid64
        // build minimal meta
        meta = tokenstorage.SessionMeta{
            UserID:    uid,
            CreatedAt: time.Now().Unix(),
            Signed:    refreshToken,
        }
    }

    // створюємо новий access token
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

    // зберігаємо нову сесію
    newMeta := tokenstorage.SessionMeta{
        UserID:    uid,
        IP:        meta.IP,
        UserAgent: meta.UserAgent,
        CreatedAt: time.Now().Unix(),
        Signed:    newRefreshSigned,
    }
    if err := tokenstorage.SaveSession(tokenstorage.RedisClient, newJTI, newMeta, auth.RefreshTokenTTL); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "can't save new refresh token"})
    }

    // видаляємо стару сесію (за oldJTI якщо є) та старий signed-key (для backward-compat)
    if oldJTI != "" {
        _ = tokenstorage.DeleteSession(tokenstorage.RedisClient, oldJTI)
    }
    _ = tokenstorage.RedisClient.Del(ctx, refreshToken).Err()

    // встановлюємо новий refresh cookie
    refreshCookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    newRefreshSigned,
        Expires:  time.Now().Add(auth.RefreshTokenTTL),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(refreshCookie)

    // встановлюємо також новий access cookie (щоб клієнтський flow з credentials: 'include' працював автоматично)
    accessCookie := &http.Cookie{
        Name:     "access_token",
        Value:    newAccess,
        Expires:  time.Now().Add(auth.AccessTokenTTL),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(accessCookie)

    // Повертаємо також access_token у тілі для backward-compat / SPA, якщо потрібно
    return c.JSON(http.StatusOK, map[string]string{"access_token": newAccess})
}

// ...existing code...
func (server *Server) Logout(c echo.Context) error {
    ctx := context.Background()

    // Видаляємо сесію за refresh token (jti або signed token)
    refreshToken := auth.ExtractToken(c)
    if refreshToken != "" {
        if jti, err := auth.ExtractJTIFromString(refreshToken); err == nil && jti != "" {
            _ = tokenstorage.DeleteSession(tokenstorage.RedisClient, jti)
        } else {
            _ = tokenstorage.RedisClient.Del(ctx, refreshToken).Err()
        }
    }

    // Очистити access_token cookie
    clearAccess := &http.Cookie{
        Name:     "access_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        Expires:  time.Unix(0, 0),
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(clearAccess)

    // Очистити refresh_token cookie
    clearRefresh := &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        Expires:  time.Unix(0, 0),
        SameSite: http.SameSiteLaxMode,
    }
    c.SetCookie(clearRefresh)

    return c.JSON(http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// 🔹 SIGN-IN логіка (перевірка користувача, створення токенів)
// Тепер SignIn повертає також jti для подальшого збереження сесії у контролері Login
func (server *Server) SignIn(email, password string) (models.Author, string, string, string, error) {
    var author models.Author
    err := server.DB.Debug().Model(models.Author{}).Where("email = ?", email).Take(&author).Error
    if err != nil {
        return author, "", "", "", err
    }

    if err := models.VerifyPassword(author.Password, password); err != nil {
        return author, "", "", "", err
    }

    // ✅ Створюємо access token
    accessToken, err := auth.CreateAccessToken(uint32(author.ID))
    if err != nil {
        return author, "", "", "", err
    }

    // ✅ створюємо refresh токен з унікальним jti
    jti := uuid.New().String()
    refreshSigned, returnedJTI, err := auth.CreateRefreshTokenWithJTI(uint32(author.ID), jti)
    if err != nil {
        return author, "", "", "", err
    }
    _ = returnedJTI // повертаємо jti нижче

    // НЕ зберігаємо signed refresh як ключ у Redis тут (ми зберігаємо session за jti у контролері Login)
    return author, accessToken, refreshSigned, jti, nil
}

// 🔹 Отримати користувача за токеном
func (server *Server) IdIfYouHaveToken(c echo.Context) error {
    // Перш за все — намагаємось дістати user id з access token (в cookie або Authorization header)
    if uid, err := auth.ExtractTokenID(c); err == nil && uid != 0 {
        author := models.Author{}
        userGotten, err := author.FindAuthorsByID(server.DB, uid)
        if err != nil {
            return c.JSON(http.StatusBadRequest, err)
        }
        return c.JSON(http.StatusOK, userGotten)
    }

    // Fallback: якщо access token відсутній/невалідний — перевіримо refresh cookie jti -> user
    refreshToken := auth.ExtractToken(c)
    if refreshToken == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing token"})
    }

    // спробуємо витягти jti
    if jti, err := auth.ExtractJTIFromString(refreshToken); err == nil && jti != "" {
        if meta, err := tokenstorage.GetSession(tokenstorage.RedisClient, jti); err == nil {
            author := models.Author{}
            userGotten, err := author.FindAuthorsByID(server.DB, uint32(meta.UserID))
            if err != nil {
                return c.JSON(http.StatusBadRequest, err)
            }
            return c.JSON(http.StatusOK, userGotten)
        }
    }

    // backward-compat: перевіримо стару схему signedToken -> userID
    ctx := context.Background()
    id, err := tokenstorage.RedisClient.Get(ctx, refreshToken).Result()
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