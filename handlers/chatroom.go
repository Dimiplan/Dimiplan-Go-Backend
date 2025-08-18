package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
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
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	return c.JSON(fiber.Map{"chatrooms": chatrooms})
}

func (h *ChatroomHandler) CreateChatroom(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	data := new(ent.Chatroom)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	chatroom, err := h.db.Chatroom.Create().SetUser(user).SetName(data.Name).Save(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	return c.JSON(fiber.Map{"chatroom_id": chatroom.ID})
}

func (h *ChatroomHandler) GetMessages(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	data := new(models.ChatroomID)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	chatroom, err := user.QueryChatrooms().Where(chatroom.ID(data.ID)).Only(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}
	messages, err := chatroom.QueryMessages().All(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	return c.JSON(fiber.Map{"messages": messages})
}

func (h *ChatroomHandler) UpdateChatroom(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	data := new(models.UpdateChatroom)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	chatroom, err := user.QueryChatrooms().Where(chatroom.ID(data.ID)).Only(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}
	if err := c.Bind().Body(chatroom); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	if _, err := chatroom.Update().SetName(data.Name).Save(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	return c.JSON(fiber.Map{"chatroom_id": chatroom.ID})
}

func (h *ChatroomHandler) RemoveChatroom(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	data := new(models.ChatroomID)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).Send(nil)
	}
	chatroom, err := user.QueryChatrooms().Where(chatroom.ID(data.ID)).Only(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Send(nil)
	}
	if err := h.db.Chatroom.DeleteOne(chatroom).Exec(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
