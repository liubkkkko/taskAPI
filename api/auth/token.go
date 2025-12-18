package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
)

// Налаштування термінів життя
const (
    AccessTokenTTL  = time.Minute * 15
    RefreshTokenTTL = time.Hour * 24 * 7
)

// CreateAccessToken — створює JWT для access (без jti)
func CreateAccessToken(userId uint32) (string, error) {
    claims := jwt.MapClaims{}
    claims["authorized"] = true
    claims["user_id"] = userId
    claims["exp"] = time.Now().Add(AccessTokenTTL).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// CreateRefreshToken — створює refresh токен з автоматичним jti
func CreateRefreshToken(userId uint32) (string, string, error) {
    jti := uuidNewString()
    return CreateRefreshTokenWithJTI(userId, jti)
}

// CreateRefreshTokenWithJTI — створює refresh токен з заданим jti
func CreateRefreshTokenWithJTI(userId uint32, jti string) (string, string, error) {
    claims := jwt.MapClaims{}
    claims["user_id"] = userId
    claims["type"] = "refresh"
    claims["jti"] = jti
    claims["exp"] = time.Now().Add(RefreshTokenTTL).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
    if err != nil {
        return "", "", err
    }
    return signed, jti, nil
}

// TokenValid перевіряє JWT (access або refresh) — повертає помилку або nil.
// Додає userID в контекст (c.Set("userID", userId)).
func TokenValid(c echo.Context) error {
    tokenString := ExtractToken(c)
    if tokenString == "" {
        log.Println("DEBUG: TokenValid - token not found")
        return fmt.Errorf("token not found")
    }

    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(os.Getenv("API_SECRET")), nil
    })
    if err != nil {
        log.Println("DEBUG: TokenValid - parse error:", err)
        return fmt.Errorf("invalid token: %v", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        log.Println("DEBUG: TokenValid - invalid claims or token")
        return errors.New("invalid or expired token")
    }

    // Отримуємо user_id
    var userId uint32
    switch v := claims["user_id"].(type) {
    case float64:
        userId = uint32(v)
    case string:
        parsed, err := strconv.ParseUint(v, 10, 32)
        if err != nil {
            log.Println("DEBUG: TokenValid - invalid user_id parse")
            return errors.New("invalid user_id claim")
        }
        userId = uint32(parsed)
    default:
        log.Println("DEBUG: TokenValid - user_id missing")
        return errors.New("user_id claim missing or invalid")
    }

    // Якщо це refresh токен — перевіряємо існування сесії в Redis по jti
    if ttype, ok := claims["type"].(string); ok && ttype == "refresh" {
        jtiRaw, ok := claims["jti"]
        if !ok {
            log.Println("DEBUG: TokenValid - refresh token missing jti")
            return errors.New("refresh token missing jti")
        }
        jti, ok := jtiRaw.(string)
        if !ok {
            log.Println("DEBUG: TokenValid - invalid jti claim")
            return errors.New("invalid jti claim")
        }
        okExists, err := tokenstorage.SessionExists(tokenstorage.RedisClient, jti)
        if err != nil {
             log.Println("DEBUG: TokenValid - redis check error:", err)
            return fmt.Errorf("redis check error: %v", err)
        }
        if !okExists {
            log.Println("DEBUG: TokenValid - refresh session not found in redis for jti:", jti)
            log.Println("Unauthorized: refresh session not found")
            return fmt.Errorf("unauthorized")
        }
    }

    c.Set("userID", userId)
    Pretty(claims)
    return nil
}

// ExtractToken — пробує дістати токен з cookie "access_token", потім Authorization, потім cookie "refresh_token", потім query param
func ExtractToken(c echo.Context) string {
    // 1. access cookie
    if cookie, err := c.Cookie("access_token"); err == nil && cookie.Value != "" {
        return cookie.Value
    }
    // 2. Authorization header
    bearer := c.Request().Header.Get("Authorization")
    if strings.HasPrefix(bearer, "Bearer ") {
        return strings.TrimPrefix(bearer, "Bearer ")
    }
    // 3. refresh cookie (fallback)
    if cookie, err := c.Cookie("refresh_token"); err == nil && cookie.Value != "" {
        return cookie.Value
    }
    // 4. query param
    if t := c.QueryParam("token"); t != "" {
        return t
    }
    return ""
}

// ExtractJTIFromString парсить JWT і повертає claim "jti".
func ExtractJTIFromString(tokenStr string) (string, error) {
    if tokenStr == "" {
        return "", errors.New("empty token")
    }
    tkn, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        // перевірка алгоритму
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        secret := os.Getenv("API_SECRET")
        return []byte(secret), nil
    })
    if err != nil {
        return "", err
    }
    claims, ok := tkn.Claims.(jwt.MapClaims)
    if !ok || !tkn.Valid {
        return "", errors.New("invalid token claims")
    }
    jtiRaw, ok := claims["jti"]
    if !ok {
        return "", errors.New("jti claim missing")
    }
    jti, ok := jtiRaw.(string)
    if !ok {
        return "", errors.New("jti claim is not a string")
    }
    return jti, nil
}

// ExtractTokenID — допоміжна функція для дістання user_id із токена
func ExtractTokenID(c echo.Context) (uint32, error) {
    tokenString := ExtractToken(c)
    if tokenString == "" {
        return 0, fmt.Errorf("token not found")
    }
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(os.Getenv("API_SECRET")), nil
    })
    if err != nil {
        return 0, err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
        if err != nil {
            return 0, err
        }
        return uint32(uid), nil
    }
    return 0, fmt.Errorf("invalid token claims")
}

// Pretty — debug для claims
func Pretty(data interface{}) {
    b, err := json.MarshalIndent(data, "", " ")
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Println(string(b))
}

// helper: uuid without importing google/uuid here to avoid compile error if not present
func uuidNewString() string {
    return uuid.New().String()
}