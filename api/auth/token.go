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
)

func CreateToken(userId uint32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires after 24 hours
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid перевіряє JWT з cookie, валідує підпис і перевіряє Redis.
// Якщо все гаразд — додає userID у echo.Context.
func TokenValid(c echo.Context) error {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return fmt.Errorf("token not found")
	}

	// Розбираємо і перевіряємо підпис токена
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

	// Отримуємо user_id з claims
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

	// Перевіряємо токен у Redis (для відкликання)
	tokenIdString := strconv.Itoa(int(userId))
	tokenExist, err := tokenstorage.CheckValueExists(tokenstorage.RedisClient, tokenIdString, tokenString)
	if err != nil {
		return fmt.Errorf("redis check error: %v", err)
	}
	if !tokenExist {
		log.Println("Unauthorized: token not found in Redis or revoked")
		return fmt.Errorf("unauthorized")
	}

	// Для зручності контролерів — кладемо userID у контекст
	c.Set("userID", userId)

	// (Не обов'язково, але залишимо для відладки)
	Pretty(claims)

	return nil
}

// ExtractToken — пробує дістати токен з cookie, потім з Authorization заголовку
func ExtractToken(c echo.Context) string {
	// cookie
	if cookie, err := c.Cookie("access_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// header: Authorization: Bearer <token>
	bearer := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}

	// optional: query param ?token=
	if t := c.QueryParam("token"); t != "" {
		return t
	}

	return ""
}

// ExtractTokenID — допоміжна функція для отримання user_id з токена (якщо потрібно окремо)
func ExtractTokenID(c echo.Context) (uint32, error) {
	tokenString := ExtractToken(c)
	if tokenString == "" {
		return 0, fmt.Errorf("token not found in cookies")
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

// Pretty виводить claims у термінал (для відладки)
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
