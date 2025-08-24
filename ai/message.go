package ai

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

func (h *AIHelper) GenerateMessage(c fiber.Ctx, model string, prompt string, data openai.ChatCompletionMessageParamUnion) (string, error) {
	if !slices.Contains(h.cfg.AIModels, model) {
		return "", fiber.NewError(fiber.StatusBadRequest, "Invalid model")
	}

	var messages []openai.ChatCompletionMessageParamUnion
	if data.OfAssistant == nil {
		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("불필요한 경우 1000 토큰 이내로 응답하세요."),
			openai.SystemMessage("LaTeX 수식은 $ 또는 $$로 감싸세요."),
			openai.UserMessage(prompt),
		}
	} else {
		messages = []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("불필요한 경우 1000 토큰 이내로 응답하세요."),
			openai.SystemMessage("LaTeX 수식은 $ 또는 $$로 감싸세요."),
			data,
			openai.UserMessage(prompt),
		}
	}

	response, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "AI response is empty")
	}

	return response.Choices[0].Message.Content, nil
}
