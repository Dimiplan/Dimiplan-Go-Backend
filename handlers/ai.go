package handlers

import (
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"dimiplan-backend/models"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

type AIHandler struct {
	db  *ent.Client
	cfg *config.Config
}

func NewAIHandler(cfg *config.Config, db *ent.Client) *AIHandler {
	return &AIHandler{
		db:  db,
		cfg: cfg,
	}
}

func (h *AIHandler) AIChat(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)
	request := new(models.AIChatRequest)
	if err := c.Bind().Body(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var room *ent.Chatroom
	if request.Room != 0 {
		var err error
		room, err = h.db.Chatroom.Query().Where(chatroom.ID(request.Room)).Only(c)
		if err != nil {
			return err
		}
	} else {
		rawresponse, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("다음 프롬프트를 요약하여 채팅방 이름을 작성하세요:\n 결과는 {\"title\": \"제목\"}의 JSON 형식으로 반환"),
				openai.UserMessage(request.Prompt),
			},
			Model: h.cfg.PreAIModel,
		})
		if err != nil || rawresponse.Choices[0].Message.Content == "" {
			return err
		}
		response := new(struct {
			Title string `json:"title"`
		})
		err = sonic.Unmarshal([]byte(rawresponse.Choices[0].Message.Content), &response)
		if err != nil {
			return err
		}
		room, err = h.db.Chatroom.Create().SetUser(user).SetName(response.Title).Save(c)
		if err != nil {
			return err
		}
	}
	if !slices.Contains(h.cfg.AIModels, request.Model) {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid model")
	}
	response, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("불필요한 경우 1000 토큰 이내로 응답하세요."),
			openai.SystemMessage("LaTeX 수식은 $ 또는 $$로 감싸세요."),
			openai.UserMessage(request.Prompt),
		},
		Model: h.cfg.PreAIModel,
	})
	if err != nil || response.Choices[0].Message.Content == "" {
		return err
	}
	_, err = h.db.Message.Create().SetChatroom(room).SetSender("user").SetMessage(request.Prompt).Save(c)
	if err != nil {
		return err
	}
	chat, err := h.db.Message.Create().SetChatroom(room).SetSender("ai").SetMessage(response.Choices[0].Message.Content).Save(c)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": chat,
		"room_id": room.ID,
	})
}
