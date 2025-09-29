package handlers

import (
	"context"
	"log"
	"strconv"

	"pets_rest/internal/database"
	"pets_rest/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type ImageHandler struct {
	imageService *services.ImageService
}

type ImageHandlerDeps struct {
	fx.In
	ImageService *services.ImageService
}

func NewImageHandler(deps ImageHandlerDeps) *ImageHandler {
	return &ImageHandler{
		imageService: deps.ImageService,
	}
}

// UploadImage завантажує зображення для конкретної сутності
func (h *ImageHandler) UploadImage(c fiber.Ctx) error {
	// Отримуємо параметри
	imageableType := c.Params("type") // "user" or "listing"
	imageableIDStr := c.Params("id")

	imageableID, err := strconv.Atoi(imageableIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Валідуємо тип
	var dbType database.ImageableType
	switch imageableType {
	case "user":
		dbType = database.ImageableTypeUser
	case "listing":
		dbType = database.ImageableTypeListing
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid imageable type. Must be 'user' or 'listing'",
		})
	}

	// Отримуємо файл з форми
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Відкриваємо файл
	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer func() {
		if closeErr := fileReader.Close(); closeErr != nil {
			log.Printf("Failed to close file reader: %v", closeErr)
		}
	}()

	// Перевіряємо, чи це основне зображення
	isPrimary := c.FormValue("is_primary") == "true"

	// Завантажуємо зображення
	image, err := h.imageService.UploadImageFile(
		context.Background(),
		dbType,
		imageableID,
		fileReader,
		file.Size,
		file.Header.Get("Content-Type"),
		file.Filename,
		isPrimary,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Повертаємо результат
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Image uploaded successfully",
		"image":   image,
	})
}

// GetImages отримує всі зображення для сутності
func (h *ImageHandler) GetImages(c fiber.Ctx) error {
	// Отримуємо параметри
	imageableType := c.Params("type")
	imageableIDStr := c.Params("id")

	imageableID, err := strconv.Atoi(imageableIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	// Валідуємо тип
	var dbType database.ImageableType
	switch imageableType {
	case "user":
		dbType = database.ImageableTypeUser
	case "listing":
		dbType = database.ImageableTypeListing
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid imageable type",
		})
	}

	// Отримуємо зображення
	images, err := h.imageService.GetImagesByEntity(dbType, imageableID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Генеруємо повні URL для відповіді (опціонально)
	response := make([]fiber.Map, len(images))
	for i, img := range images {
		fullURL, _ := h.imageService.GetImageFullURL(context.Background(), img.URL)

		response[i] = fiber.Map{
			"id":         img.ID,
			"path":       img.URL, // Шлях без домену
			"full_url":   fullURL, // Повний URL для доступу
			"filename":   img.Filename,
			"size_bytes": img.SizeBytes,
			"mime_type":  img.MimeType,
			"alt_text":   img.AltText,
			"sort_order": img.SortOrder,
			"is_primary": img.IsPrimary,
			"created_at": img.CreatedAt,
			"updated_at": img.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"images": response,
	})
}

// DeleteImage видаляє зображення
func (h *ImageHandler) DeleteImage(c fiber.Ctx) error {
	imageIDStr := c.Params("imageId")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	err = h.imageService.DeleteImage(imageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}

// SetPrimaryImage встановлює зображення як основне
func (h *ImageHandler) SetPrimaryImage(c fiber.Ctx) error {
	imageIDStr := c.Params("imageId")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	err = h.imageService.SetPrimaryImage(imageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Image set as primary successfully",
	})
}
