package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"
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
	uid := c.Locals("id").(string)

	u, err := h.db.User.Query().Where(user.ID(uid)).Only(c)
	if err != nil {
		log.Errorf("Failed to retrieve user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	return c.JSON(u)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	uid := c.Locals("id").(string)

	data := new(models.UpdateUserReq)

	if err := c.Bind().Body(data); err != nil {
		log.Errorf("Failed to bind request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	h.db.User.Update().Where(user.ID(uid)).SetName(*data.Name).Exec(c)

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
