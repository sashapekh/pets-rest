package fx

import (
	"pets_rest/internal/config"
	"pets_rest/internal/storage"

	"go.uber.org/fx"
)

// NewStorageConfig creates storage configuration from app config
func NewStorageConfig(cfg *config.Config) *storage.StorageConfig {
	return &storage.StorageConfig{
		Provider:    cfg.StorageProvider,
		Endpoint:    cfg.StorageEndpoint,
		AccessKey:   cfg.StorageAccessKey,
		SecretKey:   cfg.StorageSecretKey,
		Bucket:      cfg.StorageBucket,
		Region:      cfg.StorageRegion,
		UseSSL:      cfg.StorageUseSSL,
		Credentials: cfg.StorageCredentials,
	}
}

var StorageModule = fx.Module("storage", fx.Provide(
	NewStorageConfig,
	storage.NewStorageProvider,
))
