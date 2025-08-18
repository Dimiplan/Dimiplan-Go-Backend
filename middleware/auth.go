package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func AuthMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := session.FromContext(c).Get("id").(string)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		user, err := db.User.Query().Where(user.ID(userID)).Only(c)
		if err != nil || user == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No User",
			})
		}
		fmt.Println(user)
		c.Locals("uid", user.ID)
		c.Locals("user", user)
		return c.Next()
	}
}
