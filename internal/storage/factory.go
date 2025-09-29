package storage

import "fmt"

func NewStorageProvider(cfg *StorageConfig) (StorageProvider, error) {

	switch cfg.Provider {
	case "minio":
		return NewMinioProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", cfg.Provider)
	}
}
