package fx

import (
	"pets_rest/internal/database"

	"go.uber.org/fx"
)

var DatabaseModule = fx.Module("database", fx.Provide(database.Connect))
