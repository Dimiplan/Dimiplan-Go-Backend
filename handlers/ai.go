package handlers

import (
	"bufio"
	"dimiplan-backend/ai"
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"dimiplan-backend/models"
	"fmt"
	"log"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

type AIHandler struct {
	db       *ent.Client
	aihelper *ai.AIHelper
	cfg      *config.Config
}

func NewAIHandler(cfg *config.Config, db *ent.Client) *AIHandler {
	return &AIHandler{
		db:       db,
		aihelper: ai.NewAIHelper(cfg),
		cfg:      cfg,
	}
}

func (h *AIHandler) AIChat(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	request := rawRequest.(*models.AIChatRequest)
	user := c.Locals("user").(*ent.User)
	if !slices.Contains(h.cfg.AIModels, request.Model) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid model")
	}
	/*
		if user.ProcessingData != nil {
			return nil, fiber.NewError(fiber.StatusConflict, "User is already processing data")
		}
	*/
	value, err := sonic.MarshalString(request)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest)
	}
	user.Update().SetProcessingData(value).Save(c)
	return nil, nil
}

func (h *AIHandler) StreamAIChat(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	if user.ProcessingData == nil {
		return fiber.NewError(fiber.StatusNotAcceptable, "no data")
	}
	var request models.AIChatRequest
	err := sonic.UnmarshalString(*user.ProcessingData, &request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest)
	}

	var room *ent.Chatroom
	var data openai.ChatCompletionMessageParamUnion

	if request.Room != 0 {
		var err error
		room, err = h.db.Chatroom.Query().Where(chatroom.ID(request.Room)).Only(c)
		if err != nil {
			return err
		}
		messages, err := room.QueryMessages().All(c)
		if err != nil {
			return err
		}
		data = h.aihelper.GenerateSummary(c, messages)
	} else {
		title, err := h.aihelper.GenerateTitle(c, request.Prompt)
		if err != nil {
			return err
		}
		room, err = h.db.Chatroom.Create().SetOwner(user).SetName(title).Save(c)
		if err != nil {
			return err
		}
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache, no-transform")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	if origin := string(c.Request().Header.Peek("Origin")); origin != "" {
		c.Set("Access-Control-Allow-Origin", origin)
		c.Set("Vary", "Origin")
		c.Set("Access-Control-Allow-Credentials", "true")
	}

	return c.SendStreamWriter(func(w *bufio.Writer) {
		fmt.Fprintf(w, "Room ID: %d\n", room.ID)
		if err := w.Flush(); err != nil {
			log.Print("Client disconnected!")
			return
		}
		message := h.aihelper.GenerateMessage(w, c, request.Model, request.Prompt, data)
		if message == "" {
			return
		}

		_, err = h.db.Message.Create().
			SetChatroomID(room.ID).
			SetSender("user").
			SetMessage(request.Prompt).
			Save(c)

		_, err = h.db.Message.Create().
			SetChatroomID(room.ID).
			SetSender("ai").
			SetMessage(message).
			Save(c)

		fmt.Fprintln(w, "finish-event: message saved")
	})
}
