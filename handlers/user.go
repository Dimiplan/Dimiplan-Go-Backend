package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type UserHandler struct {
	db *ent.Client
}

func NewUserHandler(db *ent.Client) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.UpdateUserReq)

	if err := c.Bind().Body(data); err != nil {
		log.Errorf("Failed to bind request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	h.db.User.UpdateOne(user).SetName(*data.Name).Exec(c)

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
