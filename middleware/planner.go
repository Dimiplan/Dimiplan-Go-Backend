package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/planner"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func QueryPlannerMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(*ent.User)
		plannerID, err := strconv.Atoi(c.Params("planner"))
		if err != nil {
			return fiber.ErrBadRequest
		}
		planner, err := user.QueryPlanners().Where(planner.ID(plannerID)).Only(c)
		if err != nil {
			if planner == nil {
				return fiber.ErrNotFound
			}
			return err
		}
		c.Locals("planner", planner)
		return c.Next()
	}
}
