package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"

	"github.com/gofiber/fiber/v3"
	// "github.com/gofiber/fiber/v3/middleware/session"
)

type UserHandler struct {
	db *ent.Client
}

func NewUserHandler(db *ent.Client) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	userID := c.Locals("id").(string)

	user, err := h.db.User.Query().Where(user.ID(userID)).Only(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) Logout(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	h.db.User.Delete().Where(user.ID(userID)).Exec(c)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *UserHandler) Protected(c fiber.Ctx) error {
	email := c.Locals("email").(string)
	return c.JSON(fiber.Map{
		"message": "Access to protected resource",
		"email":   email,
	})
}
