package routes

import (
	"pets_rest/internal/database"
	"pets_rest/internal/handlers"

	"pets_rest/internal/config"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, db *database.DB, cfg *config.Config) {
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)

	v1 := app.Group("/api/v1")

	authHandler := handlers.NewAuthHandler(db, cfg)

	auth := v1.Group("/auth")
	auth.Get("/google/login", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)
}
