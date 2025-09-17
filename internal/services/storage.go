package services

import (
	"pets_rest/internal/config"
	"pets_rest/internal/storage"

	"go.uber.org/fx"
)

type StorageService struct {
	provider storage.StorageProvider
	config   *config.Config
}
type StorageServiceDeps struct {
	fx.In
	Config   *config.Config
	Provider storage.StorageProvider
}

func NewStorageService(deps StorageServiceDeps) (*StorageService, error) {
	return &StorageService{
		provider: deps.Provider,
		config:   deps.Config,
	}, nil
}
