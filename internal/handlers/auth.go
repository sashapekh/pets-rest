package handlers

import (
	"pets_rest/internal/auth"
	"pets_rest/internal/config"
	"pets_rest/internal/database"
	"pets_rest/internal/oauth"
	"pets_rest/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"go.uber.org/fx"
)

type AuthHandler struct {
	db             *database.DB
	cfg            *config.Config
	googleProvider *oauth.GoogleProvider
	UserService    *services.UserService
}

type AuthHandlerDeps struct {
	fx.In
	Config      *config.Config
	DB          *database.DB
	UserService *services.UserService
}

func NewAuthHandler(deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		db:             deps.DB,
		cfg:            deps.Config,
		UserService:    deps.UserService,
		googleProvider: oauth.NewGoogle(),
	}
}

func (h *AuthHandler) GoogleLogin(c fiber.Ctx) error {
	sess := session.FromContext(c)

	url, err := h.googleProvider.AuthURLWithPKCEandState(sess.Session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate auth URL",
		})
	}

	return c.Redirect().To(url)
}

func (h *AuthHandler) GoogleCallback(c fiber.Ctx) error {
	sess := session.FromContext(c)

	state := c.Query("state")
	code := c.Query("code")

	u, err := h.googleProvider.HandleCallback(c, sess.Session, state, code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.UserService.FirstOrNewUserForRegister(&u)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.cfg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"access_token": token,
		"user":         user,
	})
}
