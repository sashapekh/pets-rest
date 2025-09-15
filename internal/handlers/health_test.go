package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pets_rest/internal/database"
)

func TestHealthCheck(t *testing.T) {
	// Створити mock базу даних
	mockDB := &database.MockDB{}

	// Створити handler
	handler := NewHealthHandler(mockDB)

	// Створити Fiber app
	app := fiber.New()
	app.Get("/health", handler.HealthCheck)

	// Створити тестовий запит
	req := httptest.NewRequest("GET", "/health", http.NoBody)

	// Виконати запит
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Перевірити статус код
	assert.Equal(t, 200, resp.StatusCode)

	// Закрити response body
	_ = resp.Body.Close()
}

func TestHealthCheckWithDBError(t *testing.T) {
	// Створити mock базу даних з помилкою
	mockDB := &database.MockDBWithError{}

	// Створити handler
	handler := NewHealthHandler(mockDB)

	// Створити Fiber app
	app := fiber.New()
	app.Get("/health", handler.HealthCheck)

	// Створити тестовий запит
	req := httptest.NewRequest("GET", "/health", http.NoBody)

	// Виконати запит
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Перевірити статус код (має бути 503)
	assert.Equal(t, 503, resp.StatusCode)

	// Закрити response body
	_ = resp.Body.Close()
}
