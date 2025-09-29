package fx

import (
	"context"
	"log"
	"pets_rest/internal/config"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

var ServerModule = fx.Module("server", fx.Invoke(StartServer))

func StartServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("🚀 Server starting on port %s", cfg.Port)
				if err := app.Listen(":" + cfg.Port); err != nil {
					log.Printf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
	})
}
