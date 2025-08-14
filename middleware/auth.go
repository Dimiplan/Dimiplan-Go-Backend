package middleware

import (
	"dimiplan-backend/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(sessionService *auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := sessionService.GetIDFromSession(c)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Next()
	}
}
