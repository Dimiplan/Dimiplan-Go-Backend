package middleware

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(sessionService *auth.SessionService, db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := sessionService.GetIDFromSession(c)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		user, err := db.User.Query().Where(user.ID(userID)).First(c.Context())
		if err != nil || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No User",
			})
		}
		fmt.Println(user)
		c.Locals("uid", user.ID)
		return c.Next()
	}
}
