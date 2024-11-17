package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liubkkkko/firstAPI/api/auth"
	"github.com/liubkkkko/firstAPI/api/responses"
)

func SetMiddlewareAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := auth.TokenValid(c)
		if err != nil {
			responses.ERROR(c.Response(), http.StatusUnauthorized, errors.New("unauthorized"))
			return err
		}

		return next(c)
	}
}
