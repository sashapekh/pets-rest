package handlers

import (
	"pets_rest/internal/database"
	"time"

	"github.com/gofiber/fiber/v3"
)

type HealthHandler struct {
	db database.Interface
}

func NewHealthHandler(db database.Interface) *HealthHandler {
	return &HealthHandler{db: db}
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
