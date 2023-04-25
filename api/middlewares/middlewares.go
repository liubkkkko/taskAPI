package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/responses"
)

// func SetMiddlewareJSON(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		c.Response().Header().Set("Content-Type", "application/json")
// 		return next(c)
// 	}
// }

func SetMiddlewareAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := auth.TokenValid(c)
		if err != nil {
			responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("Unauthorized1"))
			return err
		}
		return next(c)
	}
}
