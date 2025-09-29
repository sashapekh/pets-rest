package services

import (
	"context"
	"pets_rest/internal/database"
	"pets_rest/internal/oauth"
	"time"

	"go.uber.org/fx"
)

type UserProfile struct {
	*database.User
	Avatar *string          `json:"avatar"`
	Images []database.Image `json:"images,omitempty"`
}

type UserService struct {
	userRepo     *database.UserRepository
	imageService *ImageService
}

type UserServiceDeps struct {
	fx.In
	UserRepo     *database.UserRepository
	ImageService *ImageService
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		userRepo:     deps.UserRepo,
		imageService: deps.ImageService,
	}
}

func (s *UserService) FirstOrNewUserForRegister(user *oauth.User) (*database.User, error) {
	existUser, err := s.userRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if existUser != nil {
		return existUser, nil
	}

	newUser := &database.User{
		Email:     user.Email,
		Name:      &user.Name,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetUserByID(userID int) (*database.User, error) {
	return s.userRepo.GetByID(userID)
}

// GetUserAvatar returns the primary avatar for a user
func (s *UserService) GetUserAvatar(userID int) (*database.Image, error) {
	return s.imageService.GetPrimaryImage(database.ImageableTypeUser, userID)
}

// GetUserImages returns all images for a user
func (s *UserService) GetUserImages(userID int) ([]database.Image, error) {
	return s.imageService.GetImagesByEntity(database.ImageableTypeUser, userID)
}

func (s *UserService) GetUserProfile(userID int) (*UserProfile, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Отримуємо всі зображення користувача
	images, err := s.GetUserImages(userID)
	if err != nil {
		// Якщо помилка отримання зображень, повертаємо профіль без них
		return &UserProfile{
			User:   user,
			Avatar: nil,
			Images: []database.Image{},
		}, nil
	}

	// Шукаємо основне зображення (аватар)
	var avatarURL *string
	if len(images) > 0 {
		// Спочатку шукаємо primary зображення
		for _, img := range images {
			if img.IsPrimary {
				fullURL, err := s.imageService.GetImageFullURL(context.Background(), img.URL)
				if err == nil {
					avatarURL = &fullURL
				}
				break
			}
		}

		// Якщо primary не знайдено, беремо перше зображення
		if avatarURL == nil && len(images) > 0 {
			fullURL, err := s.imageService.GetImageFullURL(context.Background(), images[0].URL)
			if err == nil {
				avatarURL = &fullURL
			}
		}
	}

	return &UserProfile{
		User:   user,
		Avatar: avatarURL,
		Images: images,
	}, nil
}
