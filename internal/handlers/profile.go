package handlers

import (
	"pets_rest/internal/auth"
	"pets_rest/internal/services"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type UserProfileHandler struct {
	userService *services.UserService
}

type UserProfileHandlerDeps struct {
	fx.In
	UserService *services.UserService
}

func NewUserProfileHandler(deps UserProfileHandlerDeps) *UserProfileHandler {
	return &UserProfileHandler{
		userService: deps.UserService,
	}
}

func (h *UserProfileHandler) GetUserProfile(c fiber.Ctx) error {

	userId, ok := auth.GetUserID(c)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	user, err := h.userService.GetUserByID(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success", "user": user})
}
