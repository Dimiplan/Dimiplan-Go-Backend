package routes

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/config"
	"dimiplan-backend/handlers"
	"dimiplan-backend/middleware"

	"github.com/bytedance/sonic"

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
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(compress.New())

	sessionService := auth.NewSessionService(cfg.RedisConfig)

	authHandler := handlers.NewAuthHandler(cfg.OAuthConfig, sessionService)
	//	userHandler := handlers.NewUserHandler(sessionService)

	app.Static("/", "dist")

	auth := app.Group("/auth")
	auth.Get("/google", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	admin := app.Group("/admin")
	admin.Use(middleware.AuthMiddleware(sessionService))

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(sessionService))
	/*
		api.Get("/profile", userHandler.GetProfile)
		api.Post("/logout", userHandler.Logout)
		api.Get("/protected", userHandler.Protected)
	*/
	return app
}
