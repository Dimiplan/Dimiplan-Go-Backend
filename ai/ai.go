package ai

import (
	"dimiplan-backend/config"
)

type AIHelper struct {
	cfg *config.Config
}

func NewAIHelper(cfg *config.Config) *AIHelper {
	return &AIHelper{
		cfg: cfg,
	}
}
