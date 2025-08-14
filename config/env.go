package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/storage/redis"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port           string
	OAuthConfig    *oauth2.Config
	RedisConfig    *redis.Config
	DatabaseString string
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
