package handlers

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/planner"
	"dimiplan-backend/models"

	"github.com/gofiber/fiber/v3"
)

type PlannerHandler struct {
	db *ent.Client
}

func NewPlannerHandler(db *ent.Client) *PlannerHandler {
	return &PlannerHandler{
		db: db,
	}
}

func (h *PlannerHandler) GetPlanner(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	// user에 Edges로 연결된 Planners 조회
	if planners, err := user.QueryPlanners().All(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve planners",
		})
	} else {
		return c.JSON(planners)
	}
}

func (h *PlannerHandler) CreatePlanner(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.CreatePlannerReq)
	if err := c.Bind().Body(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	planner, err := h.db.Planner.Create().
		SetType(data.Type).
		SetName(data.Name).
		Save(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create planner",
		})
	}

	// User에 새로운 Planner 연결
	if err := h.db.User.UpdateOne(user).AddPlanners(planner).Exec(c); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to attach planner to user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(planner)
}

func (h *PlannerHandler) GetTasks(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.GetTasksReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if planner, err := user.QueryPlanners().Where(planner.ID(data.ID)).Only(c); err != nil {
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

func (h *PlannerHandler) RenamePlanner(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.RenamePlannerReq)

	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if planner, err := user.QueryPlanners().Where(planner.ID(data.ID)).Only(c); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Planner not found",
		})
	} else {
		planner.Update().SetName(data.Name).Exec(c)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Planner renamed successfully",
		})
	}
}

func (h *PlannerHandler) DeletePlanner(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.DeletePlannerReq)
	if err := c.Bind().All(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if planner, err := user.QueryPlanners().Where(planner.ID(data.ID)).Only(c); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Planner not found",
		})
	} else {
		if err := h.db.Planner.DeleteOne(planner).Exec(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete planner",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Planner deleted successfully",
		})
	}
}
