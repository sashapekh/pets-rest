package fx

import (
	"pets_rest/internal/database"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Module("repositories", fx.Provide(
	database.NewUserRepository,
))
