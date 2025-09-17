package fx

import (
	"pets_rest/internal/services"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module("services", fx.Provide(
	services.NewStorageService,
	services.NewUserService,
))
