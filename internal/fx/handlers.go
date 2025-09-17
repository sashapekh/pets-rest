package fx

import (
	"pets_rest/internal/handlers"

	"go.uber.org/fx"
)

var HandlerModule = fx.Module("handlers", fx.Provide(
	handlers.NewHealthHandler,
	handlers.NewAuthHandler,
	handlers.NewUserProfileHandler,
))
