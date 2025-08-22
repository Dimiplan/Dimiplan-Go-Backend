package routes

import (
	"dimiplan-backend/config"
	"dimiplan-backend/ent"
	"dimiplan-backend/handlers"
	"dimiplan-backend/middleware"
	"dimiplan-backend/models"
	"dimiplan-backend/openapi"
	"fmt"
	"os"
	"strings"

	"slices"

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
	apiWrapper := openapi.NewWrapper(api)
	api.Use(middleware.AuthMiddleware(db))

	apiWrapper.Route("/user").
		Get(userHandler.GetUser, nil, ent.User{}, 200).
		Patch(userHandler.UpdateUser, new(models.UpdateUserRequest), ent.User{}, 204)

	api.Use("/ai/chatroom/:id", middleware.QueryChatroomMiddleware(db))
	apiWrapper.Route("/ai").
		Post(aiHandler.AIChat, new(models.AIChatRequest), models.AIChatResponse{}, 200).
		Route("/chatroom").
		Get(chatroomHandler.ListChatrooms, nil, models.ListChatroomsResponse{}, 200).
		Post(chatroomHandler.CreateChatroom, new(models.CreateChatroomRequest), models.CreateChatroomResponse{}, 201).
		Route("/:id").
		Get(chatroomHandler.GetMessages, new(models.GetMessagesRequest), models.GetMessagesResponse{}, 200).
		Patch(chatroomHandler.UpdateChatroom, new(models.UpdateChatroomRequest), nil, 204).
		Delete(chatroomHandler.RemoveChatroom, new(models.RemoveChatroomRequest), nil, 204)

	api.Use("/planner/:planner", middleware.QueryPlannerMiddleware(db))
	api.Use("/planner/:planner/:task", middleware.QueryTaskMiddleware(db))

	apiWrapper.Route("/planner").
		Get(plannerHandler.GetPlanners, nil, models.GetPlannersResponse{}, 200).
		Post(plannerHandler.CreatePlanner, new(models.CreatePlannerRequest), nil, 201).
		Route("/:planner").
		Get(plannerHandler.GetTasks, new(models.GetTasksRequest), models.GetTasksResponse{}, 200).
		Post(plannerHandler.CreateTask, new(models.CreateTaskRequest), nil, 201).
		Patch(plannerHandler.UpdatePlanner, new(models.RenamePlannerRequest), nil, 204).
		Delete(plannerHandler.DeletePlanner, new(models.DeletePlannerRequest), nil, 204).
		Route("/:task").
		Patch(plannerHandler.UpdateTask, new(models.UpdateTaskRequest), nil, 204).
		Delete(plannerHandler.DeleteTask, new(models.DeleteTaskRequest), nil, 204)

	os.Remove("./openapi.yaml")
	file, err := os.OpenFile("./openapi.yaml", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, apiWrapper.APIDocs())

	app.Get("/*", func(c fiber.Ctx) error {
		if slices.Contains(strings.Split(c.Path(), "/"), "api") {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendFile("dist/index.html", fiber.SendFile{
			Compress: true,
		})
	})

	return app
}
