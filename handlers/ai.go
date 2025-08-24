package handlers

import (
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"dimiplan-backend/models"
	"regexp"
	"slices"
	"strings"

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

func (h *AIHandler) AIChat(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	request := rawRequest.(*models.AIChatRequest)
	user := c.Locals("user").(*ent.User)
	var room *ent.Chatroom

	if request.Room != 0 {
		var err error
		room, err = h.db.Chatroom.Query().Where(chatroom.ID(request.Room)).Only(c)
		if err != nil {
			return nil, err
		}
	} else {
		// 채팅방 제목 생성을 위한 더 강화된 프롬프트
		rawresponse, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(`다음 사용자 프롬프트를 바탕으로 적절한 채팅방 제목을 생성해주세요.

규칙:
1. 반드시 다음 JSON 형식으로만 응답해야 합니다: {"title": "제목"}
2. 제목은 20자 이내로 작성
3. 수학 기호나 특수 문자는 피하고 일반 텍스트로만 작성
4. JSON 형식이 아닌 다른 텍스트는 절대 포함하지 마세요

예시: {"title": "수학 함수 질문"}`),
				openai.UserMessage(request.Prompt),
			},
			Model:       h.cfg.PreAIModel,
			Temperature: openai.Float(0.3), // 더 일관된 응답을 위해 temperature 낮춤
		})
		if err != nil {
			return nil, err
		}

		content := rawresponse.Choices[0].Message.Content
		if content == "" {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "AI response is empty")
		}

		// JSON 응답 정리 및 파싱
		content = h.cleanJSONResponse(content)

		response := new(struct {
			Title string `json:"title"`
		})

		err = sonic.Unmarshal([]byte(content), &response)
		if err != nil {
			// JSON 파싱 실패 시 기본 제목 사용
			response.Title = h.generateFallbackTitle(request.Prompt)
		}

		// 제목이 비어있거나 너무 긴 경우 처리
		if response.Title == "" || len(response.Title) > 50 {
			response.Title = h.generateFallbackTitle(request.Prompt)
		}

		room, err = h.db.Chatroom.Create().SetOwner(user).SetName(response.Title).Save(c)
		if err != nil {
			return nil, err
		}
	}

	// 선택된 모델이 유효한지 확인
	if !slices.Contains(h.cfg.AIModels, request.Model) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid model")
	}

	// 실제 AI 응답 생성
	response, err := h.cfg.AIClient.Chat.Completions.New(c, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("불필요한 경우 1000 토큰 이내로 응답하세요."),
			openai.SystemMessage("LaTeX 수식은 $ 또는 $$로 감싸세요."),
			openai.UserMessage(request.Prompt),
		},
		Model: request.Model, // 사용자가 선택한 모델 사용
	})
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "AI response is empty")
	}

	// 사용자 메시지 저장
	_, err = h.db.Message.Create().SetChatroom(room).SetSender("user").SetMessage(request.Prompt).Save(c)
	if err != nil {
		return nil, err
	}

	// AI 응답 저장
	chat, err := h.db.Message.Create().SetChatroom(room).SetSender("ai").SetMessage(response.Choices[0].Message.Content).Save(c)
	if err != nil {
		return nil, err
	}

	return models.AIChatResponse{
		Message: chat,
		RoomID:  room.ID,
	}, nil
}

// JSON 응답을 정리하는 헬퍼 함수
func (h *AIHandler) cleanJSONResponse(content string) string {
	// 코드 블록 마크다운 제거
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// JSON 부분만 추출하려고 시도
	jsonStart := strings.Index(content, "{")
	jsonEnd := strings.LastIndex(content, "}")

	if jsonStart >= 0 && jsonEnd >= 0 && jsonEnd > jsonStart {
		content = content[jsonStart : jsonEnd+1]
	}

	return content
}

// 대체 제목 생성 함수
func (h *AIHandler) generateFallbackTitle(prompt string) string {
	// 프롬프트에서 첫 번째 문장이나 주요 키워드 추출
	title := strings.TrimSpace(prompt)

	// 줄바꿈 제거
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")

	// 연속된 공백 제거
	re := regexp.MustCompile(`\s+`)
	title = re.ReplaceAllString(title, " ")

	// 20자로 제한
	if len(title) > 20 {
		title = title[:20] + "..."
	}

	// 빈 문자열이면 기본값 사용
	if title == "" {
		title = "새로운 채팅"
	}

	return title
}
