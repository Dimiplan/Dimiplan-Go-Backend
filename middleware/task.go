package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/task"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func QueryTaskMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		planner := c.Locals("planner").(*ent.Planner)
		taskID, err := strconv.Atoi(c.Params("task"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}
		task, err := planner.QueryTasks().Where(task.ID(taskID)).Only(c)
		if err != nil {
			log.Error(err)
			if task == nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Task not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve task",
			})
		}
		c.Locals("task", task)
		return c.Next()
	}
}
