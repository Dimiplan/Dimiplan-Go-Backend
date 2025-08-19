package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
)

type ChatroomHandler struct {
	db *ent.Client
}

func NewChatroomHandler(db *ent.Client) *ChatroomHandler {
	return &ChatroomHandler{db: db}
}

func (h *ChatroomHandler) ListChatrooms(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	chatrooms, err := h.db.User.QueryChatrooms(user).All(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"chatrooms": chatrooms})
}

func (h *ChatroomHandler) CreateChatroom(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	data := new(ent.Chatroom)
	if err := c.Bind().All(data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	chatroom, err := h.db.Chatroom.Create().SetUser(user).SetName(data.Name).Save(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"chatroom_id": chatroom.ID})
}

func (h *ChatroomHandler) GetMessages(c fiber.Ctx) error {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	messages, err := chatroom.QueryMessages().All(c)
	if err != nil {
		if messages == nil {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"messages": messages})
}

func (h *ChatroomHandler) UpdateChatroom(c fiber.Ctx) error {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	data := new(models.UpdateChatroom)
	if err := c.Bind().All(data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if _, err := chatroom.Update().SetName(data.Name).Save(c); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ChatroomHandler) RemoveChatroom(c fiber.Ctx) error {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	if err := h.db.Chatroom.DeleteOne(chatroom).Exec(c); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
