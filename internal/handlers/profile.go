package handlers

import (
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

	userProfile, err := h.userService.GetUserProfile(c.Locals("user_id").(int))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success", "user": userProfile})
}
