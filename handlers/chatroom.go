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

func (h *ChatroomHandler) ListChatrooms(request interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)
	chatrooms, err := h.db.User.QueryChatrooms(user).All(c)
	if err != nil {
		return nil, err
	}
	return models.ListChatroomsResponse{Chatrooms: chatrooms}, nil
}

func (h *ChatroomHandler) CreateChatroom(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)
	request := rawRequest.(models.CreateChatroomRequest)
	chatroom, err := h.db.Chatroom.Create().SetUser(user).SetName(request.Name).Save(c)
	if err != nil {
		return nil, err
	}
	return models.CreateChatroomResponse{ID: chatroom.ID}, nil
}

func (h *ChatroomHandler) GetMessages(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	messages, err := chatroom.QueryMessages().All(c)
	if err != nil {
		if messages == nil {
			return nil, fiber.NewError(fiber.StatusNoContent)
		}
		return nil, err
	}
	return models.GetMessagesResponse{Messages: messages}, nil
}

func (h *ChatroomHandler) UpdateChatroom(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	data := rawRequest.(models.UpdateChatroomRequest)
	if _, err := chatroom.Update().SetName(data.Name).Save(c); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *ChatroomHandler) RemoveChatroom(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	chatroom := c.Locals("chatroom").(*ent.Chatroom)
	if err := h.db.Chatroom.DeleteOne(chatroom).Exec(c); err != nil {
		return nil, err
	}
	return nil, nil
}
