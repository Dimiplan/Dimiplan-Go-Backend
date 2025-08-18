package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/planner"
	"dimiplan-backend/ent/task"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
)

func (h *PlannerHandler) GetTasks(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.GetTasksReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if planner, err := user.QueryPlanners().Where(planner.ID(data.PlannerID)).Only(c); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Planner not found",
		})
	} else {
		tasks, err := planner.QueryTasks().All(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve tasks",
			})
		}
		return c.JSON(tasks)
	}
}

func (h *PlannerHandler) CreateTask(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.CreateTaskReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if planner, err := user.QueryPlanners().Where(planner.ID(data.PlannerID)).Only(c); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Planner not found",
		})
	} else {
		task, err := h.db.Task.Create().
			SetTitle(data.Title).
			SetPriority(data.Priority).
			SetPlanner(planner).
			Save(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create task",
			})
		}
		return c.Status(fiber.StatusCreated).JSON(task)
	}
}

func (h *PlannerHandler) UpdateTask(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.UpdateTaskReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}
	planner, perr := user.QueryPlanners().Where(planner.ID(data.PlannerID)).Only(c)
	if perr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Planner not found",
		})
	}

	if task, err := planner.QueryTasks().Where(task.ID(data.TaskID)).Only(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task not found",
		})
	} else {
		task, err := task.Update().
			SetTitle(data.Title).
			SetPriority(data.Priority).
			Save(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update task",
			})
		}
		return c.Status(fiber.StatusOK).JSON(task)
	}
}

func (h *PlannerHandler) DeleteTask(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.DeleteTaskReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	planner, perr := user.QueryPlanners().Where(planner.ID(data.PlannerID)).Only(c)
	if perr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Planner not found",
		})
	}

	if task, err := planner.QueryTasks().Where(task.ID(data.TaskID)).Only(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task not found",
		})
	} else {
		if err := h.db.Task.DeleteOne(task).Exec(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete task",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Task deleted successfully",
		})
	}
}
