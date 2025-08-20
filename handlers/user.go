package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
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
		return fiber.ErrBadRequest
	}

	h.db.User.UpdateOne(user).SetName(*data.Name).Exec(c)

	return c.SendStatus(fiber.StatusNoContent)
}
