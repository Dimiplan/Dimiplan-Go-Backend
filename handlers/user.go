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

func (h *UserHandler) GetUser(request interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)

	return *user, nil
}

func (h *UserHandler) UpdateUser(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)
	request := rawRequest.(*models.UpdateUserRequest)

	h.db.User.UpdateOne(user).SetName(*request.Name).Exec(c)

	return nil, nil
}
