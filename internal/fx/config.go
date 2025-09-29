package fx

import (
	"pets_rest/internal/config"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config", fx.Provide(config.Load))
