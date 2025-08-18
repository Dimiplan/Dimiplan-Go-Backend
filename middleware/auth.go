package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/log"
)

func AuthMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		ID := session.FromContext(c).Get("id")
		var userID string
		if ID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		} else{
			userID = ID.(string)
		}
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
		log.Info(user)
		c.Locals("uid", user.ID)
		c.Locals("user", user)
		return c.Next()
	}
}
