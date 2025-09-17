package services

import (
	"pets_rest/internal/database"
	"pets_rest/internal/oauth"
	"time"

	"go.uber.org/fx"
)

type UserService struct {
	userRepo *database.UserRepository
}

type UserServiceDeps struct {
	fx.In
	UserRepo *database.UserRepository
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{userRepo: deps.UserRepo}
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
