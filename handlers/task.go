package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
)

func (h *PlannerHandler) GetTasks(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	planner := c.Locals("planner").(*ent.Planner)

	tasks, err := planner.QueryTasks().All(c)
	if err != nil {
		return nil, err
	}
	return models.GetTasksResponse{Tasks: tasks}, nil
}

func (h *PlannerHandler) CreateTask(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	planner := c.Locals("planner").(*ent.Planner)

	request := rawRequest.(*models.CreateTaskRequest)
	_, err := h.db.Task.Create().
		SetTitle(request.Title).
		SetPriority(request.Priority).
		SetPlanner(planner).
		Save(c)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *PlannerHandler) UpdateTask(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	task := c.Locals("task").(*ent.Task)
	request := rawRequest.(*models.UpdateTaskRequest)
	builder := task.Update()
	if request.Title != "" {
		builder = builder.SetTitle(request.Title)
	}
	if request.Priority != 0 {
		builder = builder.SetPriority(request.Priority)
	}
	if _, err := builder.Save(c); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *PlannerHandler) DeleteTask(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	task := c.Locals("task").(*ent.Task)
	if err := h.db.Task.DeleteOne(task).Exec(c); err != nil {
		return nil, err
	}
	return nil, nil
}
