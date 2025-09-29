package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"pets_rest/internal/database"
)

// ImageService handles operations related to images
type ImageService struct {
	db             *database.DB
	storageService *StorageService
}

// NewImageService creates a new ImageService
func NewImageService(db *database.DB, storageService *StorageService) *ImageService {
	return &ImageService{
		db:             db,
		storageService: storageService,
	}
}

// CreateImageRequest represents a request to create a new image
type CreateImageRequest struct {
	ImageableType database.ImageableType `json:"imageable_type"`
	ImageableID   int                    `json:"imageable_id"`
	URL           string                 `json:"url"`
	Filename      *string                `json:"filename,omitempty"`
	SizeBytes     *int                   `json:"size_bytes,omitempty"`
	MimeType      *string                `json:"mime_type,omitempty"`
	AltText       *string                `json:"alt_text,omitempty"`
	SortOrder     int                    `json:"sort_order"`
	IsPrimary     bool                   `json:"is_primary"`
}

// CreateImage creates a new image record in the database
func (s *ImageService) CreateImage(req CreateImageRequest) (*database.Image, error) {
	// If this is marked as primary, unset any existing primary images for this entity
	if req.IsPrimary {
		if err := s.unsetPrimaryImages(req.ImageableType, req.ImageableID); err != nil {
			return nil, fmt.Errorf("failed to unset existing primary images: %w", err)
		}
	}

	query := `
		INSERT INTO images (imageable_type, imageable_id, url, filename, size_bytes, mime_type, alt_text, sort_order, is_primary, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, imageable_type, imageable_id, url, filename, size_bytes, mime_type, alt_text, sort_order, is_primary, created_at, updated_at
	`

	var image database.Image
	err := s.db.Get(&image, query,
		req.ImageableType,
		req.ImageableID,
		req.URL,
		req.Filename,
		req.SizeBytes,
		req.MimeType,
		req.AltText,
		req.SortOrder,
		req.IsPrimary,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	return &image, nil
}

// CreateImageFromURL downloads image from URL, uploads to storage and creates a record
func (s *ImageService) CreateImageFromURL(imageableType database.ImageableType, imageableID int, imageURL string, isPrimary bool) (*database.Image, error) {
	// Download image from URL
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status: %d", resp.StatusCode)
	}

	// Extract metadata
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // default
	}

	contentLength := resp.ContentLength
	if contentLength <= 0 {
		contentLength = 1024 * 1024 // 1MB default if unknown
	}

	// Extract filename from URL
	filename := extractFilenameFromURL(imageURL)
	if filename == "" {
		filename = "image.jpg" // default
	}

	// Upload to storage (this will save to MinIO and return path without domain)
	uploadResult, err := s.storageService.UploadFile(
		context.Background(),
		resp.Body,
		contentLength,
		contentType,
		filename,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to storage: %w", err)
	}

	// Get next sort order
	sortOrder, err := s.getNextSortOrder(imageableType, imageableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sort order: %w", err)
	}

	// Create image record with path (not full URL)
	sizeBytes := int(uploadResult.Size)
	req := CreateImageRequest{
		ImageableType: imageableType,
		ImageableID:   imageableID,
		URL:           uploadResult.Path, // Save path without domain
		Filename:      &uploadResult.Filename,
		SizeBytes:     &sizeBytes,
		MimeType:      &uploadResult.ContentType,
		SortOrder:     sortOrder,
		IsPrimary:     isPrimary,
	}

	return s.CreateImage(req)
}

// GetImagesByEntity retrieves all images for a specific entity
func (s *ImageService) GetImagesByEntity(imageableType database.ImageableType, imageableID int) ([]database.Image, error) {
	query := `
		SELECT id, imageable_type, imageable_id, url, filename, size_bytes, mime_type, alt_text, sort_order, is_primary, created_at, updated_at
		FROM images
		WHERE imageable_type = $1 AND imageable_id = $2
		ORDER BY sort_order ASC, created_at ASC
	`

	var images []database.Image
	err := s.db.Select(&images, query, imageableType, imageableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	return images, nil
}

// GetPrimaryImage retrieves the primary image for a specific entity
func (s *ImageService) GetPrimaryImage(imageableType database.ImageableType, imageableID int) (*database.Image, error) {
	query := `
		SELECT id, imageable_type, imageable_id, url, filename, size_bytes, mime_type, alt_text, sort_order, is_primary, created_at, updated_at
		FROM images
		WHERE imageable_type = $1 AND imageable_id = $2 AND is_primary = TRUE
		LIMIT 1
	`

	var image database.Image
	err := s.db.Get(&image, query, imageableType, imageableID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No primary image found
		}
		return nil, fmt.Errorf("failed to get primary image: %w", err)
	}

	return &image, nil
}

// SetPrimaryImage sets an image as primary for an entity
func (s *ImageService) SetPrimaryImage(imageID int) error {
	// First get the image to know its entity
	var image database.Image
	err := s.db.Get(&image, "SELECT * FROM images WHERE id = $1", imageID)
	if err != nil {
		return fmt.Errorf("failed to find image: %w", err)
	}

	// Unset any existing primary images for this entity
	if err := s.unsetPrimaryImages(image.ImageableType, image.ImageableID); err != nil {
		return fmt.Errorf("failed to unset existing primary images: %w", err)
	}

	// Set this image as primary
	query := `UPDATE images SET is_primary = TRUE, updated_at = NOW() WHERE id = $1`
	_, err = s.db.Exec(query, imageID)
	if err != nil {
		return fmt.Errorf("failed to set primary image: %w", err)
	}

	return nil
}

// DeleteImage deletes an image by ID
func (s *ImageService) DeleteImage(imageID int) error {
	query := `DELETE FROM images WHERE id = $1`
	result, err := s.db.Exec(query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// DeleteImagesByEntity deletes all images for a specific entity
func (s *ImageService) DeleteImagesByEntity(imageableType database.ImageableType, imageableID int) error {
	query := `DELETE FROM images WHERE imageable_type = $1 AND imageable_id = $2`
	_, err := s.db.Exec(query, imageableType, imageableID)
	if err != nil {
		return fmt.Errorf("failed to delete images: %w", err)
	}

	return nil
}

// UpdateImageOrder updates the sort order of images
func (s *ImageService) UpdateImageOrder(imageableType database.ImageableType, imageableID int, imageOrders []struct {
	ID    int `json:"id"`
	Order int `json:"order"`
}) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	for _, order := range imageOrders {
		query := `
			UPDATE images 
			SET sort_order = $1, updated_at = NOW() 
			WHERE id = $2 AND imageable_type = $3 AND imageable_id = $4
		`
		_, err := tx.Exec(query, order.Order, order.ID, imageableType, imageableID)
		if err != nil {
			return fmt.Errorf("failed to update image order: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper functions

func (s *ImageService) unsetPrimaryImages(imageableType database.ImageableType, imageableID int) error {
	query := `UPDATE images SET is_primary = FALSE, updated_at = NOW() WHERE imageable_type = $1 AND imageable_id = $2 AND is_primary = TRUE`
	_, err := s.db.Exec(query, imageableType, imageableID)
	return err
}

func (s *ImageService) getNextSortOrder(imageableType database.ImageableType, imageableID int) (int, error) {
	var maxOrder sql.NullInt64
	query := `SELECT MAX(sort_order) FROM images WHERE imageable_type = $1 AND imageable_id = $2`
	err := s.db.Get(&maxOrder, query, imageableType, imageableID)
	if err != nil {
		return 0, err
	}

	if maxOrder.Valid {
		return int(maxOrder.Int64) + 1, nil
	}

	return 0, nil
}

func extractFilenameFromURL(url string) string {
	// Extract filename from URL path
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// Remove query parameters
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
		return filename
	}
	return ""
}

// GetImageFullURL повертає повний URL для зображення за його шляхом
func (s *ImageService) GetImageFullURL(ctx context.Context, imagePath string) (string, error) {
	return s.storageService.GetFileURL(ctx, imagePath)
}

// UploadImageFile завантажує файл безпосередньо та створює запис
func (s *ImageService) UploadImageFile(ctx context.Context, imageableType database.ImageableType, imageableID int, reader io.Reader, size int64, contentType, filename string, isPrimary bool) (*database.Image, error) {
	// Upload to storage
	uploadResult, err := s.storageService.UploadFile(ctx, reader, size, contentType, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image to storage: %w", err)
	}

	// Get next sort order
	sortOrder, err := s.getNextSortOrder(imageableType, imageableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sort order: %w", err)
	}

	// Create image record with path (not full URL)
	sizeBytes := int(uploadResult.Size)
	req := CreateImageRequest{
		ImageableType: imageableType,
		ImageableID:   imageableID,
		URL:           uploadResult.Path, // Save path without domain
		Filename:      &uploadResult.Filename,
		SizeBytes:     &sizeBytes,
		MimeType:      &uploadResult.ContentType,
		SortOrder:     sortOrder,
		IsPrimary:     isPrimary,
	}

	return s.CreateImage(req)
}
