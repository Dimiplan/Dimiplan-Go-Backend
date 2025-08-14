package routes

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/handlers"
	"dimiplan-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, cfg *config.Config, db *ent.Client) *fiber.App {
	sessionService := auth.NewSessionService(cfg.RedisConfig)

	authHandler := handlers.NewAuthHandler(cfg.OAuthConfig, sessionService, db)
	userHandler := handlers.NewUserHandler(sessionService, db)

	auth := app.Group("/auth")
	auth.Get("/google", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	admin := app.Group("/admin")
	admin.Use(middleware.AuthMiddleware(sessionService, db))

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(sessionService, db))
	api.Get("/profile", userHandler.GetProfile)
	api.Post("/logout", userHandler.Logout)
	api.Get("/protected", userHandler.Protected)
	return app
}
