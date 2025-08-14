package config

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Port        string
	JWTSecret   []byte
	RedisClient *redis.Client
	OAuthConfig *oauth2.Config
}

func Load() *Config {
	godotenv.Load()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

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

	return &Config{
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   []byte(getEnv("JWT_SECRET", "your-secret-key")),
		RedisClient: redisClient,
		OAuthConfig: oauthConfig,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
