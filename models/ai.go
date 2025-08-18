package models

type AIChatRequest struct {
	Prompt string `json:"prompt"`
	Room   int    `json:"room"`
	Model  string `json:"model"`
	Search bool   `json:"search"`
}
