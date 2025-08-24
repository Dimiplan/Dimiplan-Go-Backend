package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/storage/redis/v3"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port           string
	OAuthConfig    *oauth2.Config
	RedisConfig    *redis.Config
	DatabaseString string
	AIClient       openai.Client
	PreAIModel     string
	AIModels       []string
}

func Load() *Config {
	godotenv.Load()

	oauthConfig := &oauth2.Config{
		ClientID:     getEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  getEnv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	redisConfig := &redis.Config{
		Host:     getEnv("REDIS_HOST"),
		Port:     getEnvAsInt("REDIS_PORT"),
		Password: "",
	}

	return &Config{
		Port:        getEnv("PORT"),
		OAuthConfig: oauthConfig,
		RedisConfig: redisConfig,
		DatabaseString: fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable",
			getEnv("DB_USER"),
			getEnv("DB_HOST"),
			getEnvAsInt("DB_PORT"),
			getEnv("DB_NAME"),
		),
		AIClient: openai.NewClient(
			option.WithBaseURL("https://openrouter.ai/api/v1"),
			option.WithAPIKey(getEnv("OPENROUTER_API_KEY")),
		),
		PreAIModel: "openai/gpt-oss-120b",
		AIModels: []string{
			"anthropic/claude-3.5-haiku",
			"deepseek/deepseek-chat-v3.1",
			"deepseek/deepseek-r1-0528",
			"google/gemini-2.5-flash",
			"meta-llama/llama-4-maverick",
			"microsoft/phi-4-reasoning-plus",
			"mistralai/devstral-medium",
			"mistralai/magistral-medium-2506:thinking",
			"moonshotai/kimi-k2",
			"openai/gpt-5-chat",
			"openai/gpt-oss-120b",
			"qwen/qwen3-coder",
			"x-ai/grok-3-mini",
		},
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s not found", key))
}

func getEnvAsInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	panic(fmt.Sprintf("Environment variable %s not found", key))
}
