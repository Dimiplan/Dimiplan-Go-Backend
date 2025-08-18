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
		ClientID:     getEnv("GOOGLE_CLIENT_ID", "your-client-id"),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "your-client-secret"),
		RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	redisConfig := &redis.Config{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		OAuthConfig: oauthConfig,
		RedisConfig: redisConfig,
		DatabaseString: fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnvAsInt("DB_PORT", 5432),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_NAME", "dimiplan"),
		),
		AIClient: openai.NewClient(
			option.WithBaseURL("https://openrouter.ai/api/v1"),
			option.WithAPIKey(getEnv("OPENAI_API_KEY", "your-api-key")),
		),
		PreAIModel: "openai/gpt-oss-120b",
		AIModels: []string{
			"anthropic/claude-3.5-haiku",
			"deepseek/deepseek-prover-v2",
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
