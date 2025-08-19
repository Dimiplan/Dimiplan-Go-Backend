package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func (h *PlannerHandler) GetTasks(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	tasks, err := planner.QueryTasks().All(c)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tasks",
		})
	}
	return c.JSON(tasks)
}

func (h *PlannerHandler) CreateTask(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	data := new(models.CreateTaskReq)
	if err := c.Bind().All(data); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}
	_, err := h.db.Task.Create().
		SetTitle(data.Title).
		SetPriority(data.Priority).
		SetPlanner(planner).
		Save(c)

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task",
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *PlannerHandler) UpdateTask(c fiber.Ctx) error {
	task := c.Locals("task").(*ent.Task)
	data := new(models.UpdateTaskReq)
	if err := c.Bind().All(data); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	_, err := task.Update().
		SetTitle(data.Title).
		SetPriority(data.Priority).
		Save(c)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PlannerHandler) DeleteTask(c fiber.Ctx) error {
	task := c.Locals("task").(*ent.Task)
	if err := h.db.Task.DeleteOne(task).Exec(c); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
