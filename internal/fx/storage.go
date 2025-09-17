package fx

import (
	"pets_rest/internal/storage"

	"go.uber.org/fx"
)

var StorageModule = fx.Module("storage", fx.Provide(storage.NewStorageProvider))
