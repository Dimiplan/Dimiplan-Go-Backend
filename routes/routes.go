package routes

import (
	"context"
	"dimiplan-backend/auth"
	"dimiplan-backend/config"
	"dimiplan-backend/handlers"
	"dimiplan-backend/middleware"
	"dimiplan-backend/storage"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(compress.New())

	ctx := context.Background()
	if err := cfg.RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	oauthService := auth.NewOAuthService(cfg.OAuthConfig, cfg.RedisClient)
	storageService := storage.NewRedisService(cfg.RedisClient)

	authHandler := handlers.NewAuthHandler(oauthService, jwtService, storageService)
	userHandler := handlers.NewUserHandler(storageService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Fiber Google OAuth JWT Redis System",
			"endpoints": fiber.Map{
				"login":    "/auth/google",
				"callback": "/auth/google/callback",
				"profile":  "/api/profile",
				"logout":   "/api/logout",
			},
		})
	})

	app.Get("/auth/google", authHandler.GoogleLogin)
	app.Get("/auth/google/callback", authHandler.GoogleCallback)

	api := app.Group("/api")
	api.Use(middleware.JWT(jwtService))

	api.Get("/profile", userHandler.GetProfile)
	api.Post("/logout", userHandler.Logout)
	api.Get("/protected", userHandler.Protected)

	return app
}
