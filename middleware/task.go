package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/task"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func QueryTaskMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		planner := c.Locals("planner").(*ent.Planner)
		taskID, err := strconv.Atoi(c.Params("task"))
		if err != nil {
			return fiber.ErrBadRequest
		}
		task, err := planner.QueryTasks().Where(task.ID(taskID)).Only(c)
		if err != nil {
			if task == nil {
				return fiber.ErrNotFound
			}
			return err
		}
		c.Locals("task", task)
		return c.Next()
	}
}
