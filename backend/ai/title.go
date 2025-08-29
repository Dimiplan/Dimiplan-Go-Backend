package ai

import (
	"regexp"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

func (h *AIHelper) GenerateTitle(c fiber.Ctx, prompt string) (string, error) {
	rawresponse, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`다음 사용자 프롬프트를 바탕으로 적절한 채팅방 제목을 생성해주세요.

규칙:
1. 반드시 다음 JSON 형식으로만 응답해야 합니다: {"title": "제목"}
2. 제목은 20자 이내로 작성
3. 수학 기호나 특수 문자는 피하고 일반 텍스트로만 작성
4. JSON 형식이 아닌 다른 텍스트는 절대 포함하지 마세요

예시: {"title": "수학 함수 질문"}`),
			openai.UserMessage(prompt),
		},
		Model:       h.cfg.PreAIModel,
		Temperature: openai.Float(0.3),
	})
	if err != nil {
		return "", err
	}

	content := rawresponse.Choices[0].Message.Content
	if content == "" {
		return "", fiber.NewError(fiber.StatusInternalServerError, "AI response is empty")
	}

	content = h.cleanJSONResponse(content)

	response := new(struct {
		Title string `json:"title"`
	})

	err = sonic.Unmarshal([]byte(content), &response)
	if err != nil {
		response.Title = h.generateFallbackTitle(prompt)
	}

	if response.Title == "" || len(response.Title) > 50 {
		response.Title = h.generateFallbackTitle(prompt)
	}
	return response.Title, nil
}

func (h *AIHelper) cleanJSONResponse(content string) string {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	jsonStart := strings.Index(content, "{")
	jsonEnd := strings.LastIndex(content, "}")

	if jsonStart >= 0 && jsonEnd >= 0 && jsonEnd > jsonStart {
		content = content[jsonStart : jsonEnd+1]
	}

	return content
}

func (h *AIHelper) generateFallbackTitle(prompt string) string {
	title := strings.TrimSpace(prompt)

	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")

	re := regexp.MustCompile(`\s+`)
	title = re.ReplaceAllString(title, " ")

	if len(title) > 20 {
		title = title[:20] + "..."
	}

	if title == "" {
		title = "새로운 채팅"
	}

	return title
}
