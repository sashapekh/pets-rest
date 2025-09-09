package handlers

import (
	"pets_rest/internal/config"
	"pets_rest/internal/database"
	"pets_rest/internal/oauth"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type AuthHandler struct {
	db             *database.DB
	cfg            *config.Config
	googleProvider *oauth.GoogleProvider
}

func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db: db, cfg: cfg,
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

	return c.JSON(fiber.Map{
		"user": u,
	})
}
