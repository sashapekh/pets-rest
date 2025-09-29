package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

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

type UploadResult struct {
	Path        string // Шлях без домену, наприклад: "/pets-photos/image123.jpg"
	FullURL     string // Повний URL для доступу
	Size        int64
	ContentType string
	Filename    string
}

func NewStorageService(deps StorageServiceDeps) (*StorageService, error) {
	return &StorageService{
		provider: deps.Provider,
		config:   deps.Config,
	}, nil
}

// UploadFile завантажує файл в MinIO та повертає результат з шляхом без домену
func (s *StorageService) UploadFile(ctx context.Context, reader io.Reader, size int64, contentType, originalFilename string) (*UploadResult, error) {
	// Генеруємо унікальне ім'я файлу
	key := s.generateFileKey(originalFilename)

	// Завантажуємо файл в MinIO
	err := s.provider.UploadFile(ctx, key, reader, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Генеруємо шлях без домену
	path := fmt.Sprintf("/%s/%s", s.getBucketName(), key)

	// Генеруємо повний URL (для відповіді клієнту або логування)
	fullURL, err := s.provider.GetFileURL(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get file URL: %w", err)
	}

	return &UploadResult{
		Path:        path,
		FullURL:     fullURL,
		Size:        size,
		ContentType: contentType,
		Filename:    originalFilename,
	}, nil
}

// GetFileURL повертає повний URL для файлу за шляхом
func (s *StorageService) GetFileURL(ctx context.Context, path string) (string, error) {
	// Видаляємо префікс bucket з шляху, щоб отримати key
	key := s.pathToKey(path)
	return s.provider.GetFileURL(ctx, key)
}

// DeleteFile видаляє файл за шляхом
func (s *StorageService) DeleteFile(ctx context.Context, path string) error {
	key := s.pathToKey(path)
	return s.provider.DeleteFile(ctx, key)
}

// generateFileKey генерує унікальне ім'я файлу
func (s *StorageService) generateFileKey(originalFilename string) string {
	// Отримуємо розширення файлу
	ext := filepath.Ext(originalFilename)

	// Генеруємо унікальне ім'я на базі timestamp
	timestamp := time.Now().UnixNano()

	// Очищуємо оригінальне ім'я від небезпечних символів
	baseName := strings.TrimSuffix(filepath.Base(originalFilename), ext)
	safeName := strings.ReplaceAll(baseName, " ", "_")

	return fmt.Sprintf("%d_%s%s", timestamp, safeName, ext)
}

// pathToKey конвертує шлях (наприклад "/pets-photos/file.jpg") в key ("file.jpg")
func (s *StorageService) pathToKey(path string) string {
	// Видаляємо початковий слеш та префікс bucket
	trimmed := strings.TrimPrefix(path, "/")
	bucketPrefix := s.getBucketName() + "/"
	return strings.TrimPrefix(trimmed, bucketPrefix)
}

// getBucketName повертає ім'я bucket (можна винести в конфіг)
func (s *StorageService) getBucketName() string {
	// TODO: Винести в конфігурацію
	return "pets-photos"
}
