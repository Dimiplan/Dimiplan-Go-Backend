package routes

import (
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/handlers"
	"dimiplan-backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, cfg *config.Config, db *ent.Client) *fiber.App {
	authHandler := handlers.NewAuthHandler(cfg.OAuthConfig, db)
	userHandler := handlers.NewUserHandler(db)
	aiHandler := handlers.NewAIHandler(cfg, db)
	chatroomHandler := handlers.NewChatroomHandler(db)
	plannerHandler := handlers.NewPlannerHandler(db)

	auth := app.Group("/auth")
	auth.Get("/login", authHandler.Login)
	auth.Get("/callback", authHandler.Callback)
	auth.Get("/logout", authHandler.Logout)

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(db))

	api.Route("/user").
		Get(userHandler.GetUser).
		Patch(userHandler.UpdateUser)

	api.Route("/ai").
		Post(aiHandler.AIChat).
		Route("/chatroom").
			Get(chatroomHandler.ListChatrooms).
			Post(chatroomHandler.CreateChatroom).
			Route("/:id").
				Get(chatroomHandler.GetMessages).
				Patch(chatroomHandler.UpdateChatroom).
				Delete(chatroomHandler.RemoveChatroom)

	api.Route("/planner").
		Get(plannerHandler.GetPlanners).
		Post(plannerHandler.CreatePlanner).
		Route("/:planner").
			Get(plannerHandler.GetTasks).
			Post(plannerHandler.CreateTask).
			Patch(plannerHandler.UpdatePlanner).
			Delete(plannerHandler.DeletePlanner).
			Route("/:task").
				Patch(plannerHandler.UpdateTask).
				Delete(plannerHandler.DeleteTask)

	return app
}
