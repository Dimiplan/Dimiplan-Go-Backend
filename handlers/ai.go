package handlers

import (
	"dimiplan-backend/ai"
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

type AIHandler struct {
	db       *ent.Client
	aihelper *ai.AIHelper
}

func NewAIHandler(cfg *config.Config, db *ent.Client) *AIHandler {
	return &AIHandler{
		db:       db,
		aihelper: ai.NewAIHelper(cfg),
	}
}

func (h *AIHandler) AIChat(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	request := rawRequest.(*models.AIChatRequest)
	user := c.Locals("user").(*ent.User)
	var room *ent.Chatroom
	var data openai.ChatCompletionMessageParamUnion

	if request.Room != 0 {
		var err error
		room, err = h.db.Chatroom.Query().Where(chatroom.ID(request.Room)).Only(c)
		if err != nil {
			return nil, err
		}
		messages, err := room.QueryMessages().All(c)
		if err != nil {
			return nil, err
		}
		data = h.aihelper.GenerateSummary(c, messages)
	} else {
		title, err := h.aihelper.GenerateTitle(c, request.Prompt)
		if err != nil {
			return nil, err
		}
		room, err = h.db.Chatroom.Create().SetOwner(user).SetName(title).Save(c)
		if err != nil {
			return nil, err
		}
	}

	message, err := h.aihelper.GenerateMessage(c, request.Model, request.Prompt, data)
	if err != nil {
		return nil, err
	}

	_, err = h.db.Message.Create().SetChatroom(room).SetSender("user").SetMessage(request.Prompt).Save(c)
	if err != nil {
		return nil, err
	}

	chat, err := h.db.Message.Create().SetChatroom(room).SetSender("ai").SetMessage(message).Save(c)
	if err != nil {
		return nil, err
	}

	return models.AIChatResponse{
		Message: chat,
		RoomID:  room.ID,
	}, nil
}
