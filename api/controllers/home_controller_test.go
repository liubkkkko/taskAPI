package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	// Створюємо тестовий сервер
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Створюємо об'єкт сервера та викликаємо обробник Home
	server := &Server{} 
	err := server.Home(c)

	// Перевіряємо, чи немає помилки, і що відповідь містить очікуваний текст
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "\"Welcome To This Awesome API\"\n", rec.Body.String())
}

func TestTestRout(t *testing.T) {
	// Створюємо тестовий сервер
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Створюємо об'єкт сервера та викликаємо обробник TestRout
	server := &Server{}
	err := server.TestRout(c)

	// Перевіряємо, чи немає помилки, і що відповідь містить очікуваний текст
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "\"created new routs\"\n", rec.Body.String())
}
