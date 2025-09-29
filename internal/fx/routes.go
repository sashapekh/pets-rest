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
	ImageHandler       *handlers.ImageHandler
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

	// Image routes
	images := v1.Group("/images")
	images.Use(auth.JWTMiddleware(deps.Config)) // Потребує аутентифікації

	// Upload image: POST /api/v1/images/:type/:id
	// Example: POST /api/v1/images/user/123 or POST /api/v1/images/listing/456
	images.Post("/:type/:id", deps.ImageHandler.UploadImage)

	// Get images: GET /api/v1/images/:type/:id
	images.Get("/:type/:id", deps.ImageHandler.GetImages)

	// Delete specific image: DELETE /api/v1/images/:imageId
	images.Delete("/:imageId", deps.ImageHandler.DeleteImage)

	// Set image as primary: PUT /api/v1/images/:imageId/primary
	images.Put("/:imageId/primary", deps.ImageHandler.SetPrimaryImage)
}
