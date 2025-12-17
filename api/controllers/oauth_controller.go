package controllers

import (
	"context"
	"encoding/json"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"os"

	"github.com/google/uuid"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/models"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
)

// getOAuthConfig будує конфіг для google oauth2 з env
func getOAuthConfig() *oauth2.Config {
    return &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"), // наприклад https://localhost:443/auth/google/callback
        Scopes:       []string{"openid", "email", "profile"},
        Endpoint:     google.Endpoint,
    }
}

// GoogleAuthRedirect — редіректить користувача на Google consent page
func (server *Server) GoogleAuthRedirect(c echo.Context) error {
    conf := getOAuthConfig()
    // state можна створити реальніший (збереження в redis) — тут простий uuid
    state := uuid.New().String()
    // опціонально зберегти state в redis для перевірки в callback
    // redirect to google
    url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
    return c.Redirect(http.StatusFound, url)
}

// GoogleAuthCallback — обробляє callback від Google
func (server *Server) GoogleAuthCallback(c echo.Context) error {
    conf := getOAuthConfig()
    ctx := context.Background()

    code := c.QueryParam("code")
    if code == "" {
        return c.String(http.StatusBadRequest, "missing code")
    }

    // обмін коду на токен
    tok, err := conf.Exchange(ctx, code)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token exchange failed"})
    }

    // отримати userinfo
    client := conf.Client(ctx, tok)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "userinfo fetch failed"})
    }
    defer resp.Body.Close()

    var gi struct {
        Sub           string `json:"sub"`
        Email         string `json:"email"`
        VerifiedEmail bool   `json:"email_verified"`
        Name          string `json:"name"`
        GivenName     string `json:"given_name"`
        FamilyName    string `json:"family_name"`
        Picture       string `json:"picture"`
        Locale        string `json:"locale"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&gi); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid userinfo"})
    }

    // знайти або створити автора в БД по email
    var author models.Author
    if err := server.DB.Where("email = ?", gi.Email).First(&author).Error; err != nil {
        // create new author
        author = models.Author{
            Nickname: gi.Name,
            Email:    gi.Email,
            Password: "", // пустий, OAuth-user cannot login with password
        }
        // якщо у моделі є Prepare/Validate/Save — можна викликати, але використовуємо просте створення
        if err := server.DB.Create(&author).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot create user"})
        }
    }

    // створюємо access + refresh token та зберігаємо сесію (аналог Login)
    accessToken, err := auth.CreateAccessToken(uint32(author.ID))
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot create access token"})
    }
    jti := uuid.New().String()
    refreshSigned, _, err := auth.CreateRefreshTokenWithJTI(uint32(author.ID), jti)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot create refresh token"})
    }

    // зберігаємо сесію в redis
    meta := tokenstorage.SessionMeta{
        UserID:    int(author.ID),
        IP:        c.RealIP(),
        UserAgent: c.Request().Header.Get("User-Agent"),
        CreatedAt: time.Now().Unix(),
        Signed:    refreshSigned,
    }
    
if err := tokenstorage.SaveSession(tokenstorage.RedisClient, jti, meta, auth.RefreshTokenTTL); err != nil {
    server.Logger.Warn(
        "failed to save oauth session",
        zap.Error(err),
        zap.String("jti", jti),
        zap.Int("user_id", int(author.ID)),
    )
}


    // ставимо cookie як у Login
    accessCookie := &http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Expires:  time.Now().Add(auth.AccessTokenTTL),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteNoneMode,
    }
    c.SetCookie(accessCookie)

    refreshCookie := &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshSigned,
        Expires:  time.Now().Add(auth.RefreshTokenTTL),
        HttpOnly: true,
        Secure:   false,
        Path:     "/",
        SameSite: http.SameSiteNoneMode,
    }
    c.SetCookie(refreshCookie)

    // редірект на фронтенд (приклад на /workspaces). Не передаємо токени у query.
    frontend := os.Getenv("FRONTEND_URL")
    if frontend == "" {
        frontend = "https://localhost:3000"
    }
    return c.Redirect(http.StatusFound, frontend+"/workspaces")
}