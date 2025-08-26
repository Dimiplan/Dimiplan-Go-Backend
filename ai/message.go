package ai

import (
	"bufio"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/openai/openai-go/v2"
)

func (h *AIHelper) GenerateMessage(w *bufio.Writer, c fiber.Ctx, model string, prompt string, data openai.ChatCompletionMessageParamUnion) string {
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

	stream := h.cfg.AIClient.Chat.Completions.NewStreaming(c, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})

	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		if _, ok := acc.JustFinishedContent(); ok {
			fmt.Fprintln(w)
			fmt.Fprintln(w, "finish-event: Content stream finished")
			if err := w.Flush(); err != nil {
				log.Print("Client disconnected!")
			}
			break
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Fprint(w, chunk.Choices[0].Delta.Content)
			if err := w.Flush(); err != nil {
				log.Print("Client disconnected!")
				stream.Close()
				return ""
			}
		}
	}
	return acc.Choices[0].Message.Content
}
