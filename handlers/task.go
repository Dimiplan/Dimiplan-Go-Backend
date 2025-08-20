package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
)

func (h *PlannerHandler) GetTasks(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	tasks, err := planner.QueryTasks().All(c)
	if err != nil {
		return err
	}
	return c.JSON(tasks)
}

func (h *PlannerHandler) CreateTask(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	data := new(models.CreateTaskReq)
	if err := c.Bind().All(data); err != nil {
		return fiber.ErrBadRequest
	}
	_, err := h.db.Task.Create().
		SetTitle(data.Title).
		SetPriority(data.Priority).
		SetPlanner(planner).
		Save(c)

	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *PlannerHandler) UpdateTask(c fiber.Ctx) error {
	task := c.Locals("task").(*ent.Task)
	data := new(models.UpdateTaskReq)
	if err := c.Bind().All(data); err != nil {
		return fiber.ErrBadRequest
	}

	builder := task.Update()
	if data.Title != "" {
		builder = builder.SetTitle(data.Title)
	}
	if data.Priority != 0 {
		builder = builder.SetPriority(data.Priority)
	}
	if _, err := builder.Save(c); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PlannerHandler) DeleteTask(c fiber.Ctx) error {
	task := c.Locals("task").(*ent.Task)
	if err := h.db.Task.DeleteOne(task).Exec(c); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
