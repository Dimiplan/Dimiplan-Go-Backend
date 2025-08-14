package handlers

import (
	"dimiplan-backend/storage"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	storage *storage.RedisService
}

func NewUserHandler(storage *storage.RedisService) *UserHandler {
	return &UserHandler{storage: storage}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.storage.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	h.storage.DeleteUser(userID)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *UserHandler) Protected(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	return c.JSON(fiber.Map{
		"message": "Access to protected resource",
		"email":   email,
	})
}
