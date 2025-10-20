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

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/tokenstorage"
	"github.com/google/uuid"
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
	jti := uuid.New().String()
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
		return fmt.Errorf("token not found")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
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
			return errors.New("invalid user_id claim")
		}
		userId = uint32(parsed)
	default:
		return errors.New("user_id claim missing or invalid")
	}

	// Перевіряємо, чи цей токен існує в Redis (тобто не відкликаний)
	tokenIdString := strconv.Itoa(int(userId))
	exists, err := tokenstorage.CheckValueExists(tokenstorage.RedisClient, tokenIdString, tokenString)
	if err != nil {
		return fmt.Errorf("redis check error: %v", err)
	}
	if !exists {
		log.Println("Unauthorized: token not found or revoked")
		return fmt.Errorf("unauthorized")
	}

	c.Set("userID", userId)
	Pretty(claims)
	return nil
}

// ExtractToken — пробує дістати токен з cookie "refresh_token", потім з Authorization, потім query param
func ExtractToken(c echo.Context) string {
	if cookie, err := c.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	bearer := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}
	if t := c.QueryParam("token"); t != "" {
		return t
	}
	return ""
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
