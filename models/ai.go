package models

import "dimiplan-backend/ent"

type AIChatRequest struct {
	Prompt string `json:"prompt"`
	Room   int    `json:"room"`
	Model  string `json:"model"`
	Search bool   `json:"search"`
}

type AIChatResponse struct {
	Message *ent.Message `json:"message"`
	RoomID  int          `json:"room_id"`
}
