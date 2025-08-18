package handlers

import (
	"dimiplan-backend/ent"
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

func (h* PlannerHandler) 