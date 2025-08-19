package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/planner"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func QueryPlannerMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(*ent.User)
		plannerID, err := strconv.Atoi(c.Params("planner"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid planner ID",
			})
		}
		planner, err := user.QueryPlanners().Where(planner.ID(plannerID)).Only(c)
		if err != nil {
			log.Error(err)
			if planner == nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Planner not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve planner",
			})
		}
		c.Locals("planner", planner)
		return c.Next()
	}
}
