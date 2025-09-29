package auth

import (
	"pets_rest/internal/config"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func JWTMiddleware(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		token := tokenParts[1]
		claims, err := ValidateToken(token, cfg)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		c.Locals("user_id", claims.UserID)

		return c.Next()
	}
}

func GetUserID(c fiber.Ctx) (int, bool) {
	userID, ok := c.Locals("user_id").(int)
	return userID, ok
}

func GetClaims(c fiber.Ctx) (*Claims, bool) {
	claims, ok := c.Locals("claims").(*Claims)
	return claims, ok
}
