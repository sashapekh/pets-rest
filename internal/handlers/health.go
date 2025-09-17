package handlers

import (
	"time"

	"pets_rest/internal/database"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type HealthHandler struct {
	db *database.DB
}

type HealthHandlerDeps struct {
	fx.In
	DB *database.DB
}

func NewHealthHandler(deps HealthHandlerDeps) *HealthHandler {
	return &HealthHandler{db: deps.DB}
}

func (h *HealthHandler) HealthCheck(c fiber.Ctx) error {
	// Check database health
	if err := h.db.Health(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Service is healthy",
		"time":    time.Now().Format(time.RFC3339),
	})
}
