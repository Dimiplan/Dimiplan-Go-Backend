package ai

import (
	"dimiplan-backend/ent"

	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

func (h *AIHelper) GenerateSummary(c fiber.Ctx, m []*ent.Message) openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion

	for _, msg := range m {
		messages = append(messages, openai.UserMessage(msg.Message))
	}

	response, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: append([]openai.ChatCompletionMessageParamUnion{openai.SystemMessage("메시지를 의미를 잃지 않는 선에서 최대한 요약하세요.")}, messages...),
		Model:    h.cfg.PreAIModel,
	})
	if err != nil {
		return openai.ChatCompletionMessageParamUnion{}
	}

	return openai.AssistantMessage(response.Choices[0].Message.Content)
}
