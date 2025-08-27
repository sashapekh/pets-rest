package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	// Health check
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Pets Search API is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/magic-link", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Magic link endpoint - not implemented yet"})
	})
	auth.Post("/magic-link/verify", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Magic link verify endpoint - not implemented yet"})
	})

	// Listings routes
	listings := api.Group("/listings")
	listings.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get listings - not implemented yet"})
	})
	listings.Post("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Create listing - not implemented yet"})
	})
	listings.Get("/:id", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get listing by ID - not implemented yet"})
	})
	listings.Put("/:id", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Update listing - not implemented yet"})
	})
	listings.Delete("/:id", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Delete listing - not implemented yet"})
	})

	// Public pages
	app.Get("/p/:slug", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Public page - not implemented yet"})
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
