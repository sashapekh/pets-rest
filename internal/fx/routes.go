package fx

import (
	"fmt"
	"pets_rest/internal/auth"
	"pets_rest/internal/config"
	"pets_rest/internal/handlers"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

var RouteModule = fx.Module("routes", fx.Invoke(SetupRoutes))

type RouteDependencies struct {
	fx.In
	HealthHandler      *handlers.HealthHandler
	AuthHandler        *handlers.AuthHandler
	UserProfileHandler *handlers.UserProfileHandler
	Config             *config.Config
}

func SetupRoutes(app *fiber.App, deps RouteDependencies) {
	fmt.Println("SetupRoutes called by FX")

	// Health check route
	app.Get("/health", deps.HealthHandler.HealthCheck)

	// API v1 routes
	v1 := app.Group("/api/v1")

	// Auth routes
	authRouter := v1.Group("/auth")
	authRouter.Get("/google/login", deps.AuthHandler.GoogleLogin)
	authRouter.Get("/google/callback", deps.AuthHandler.GoogleCallback)

	// User profile routes
	profile := v1.Group("/profile")
	profile.Use(auth.JWTMiddleware(deps.Config))
	profile.Get("/", deps.UserProfileHandler.GetUserProfile)
}
