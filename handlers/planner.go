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

func (h *PlannerHandler) GetPlanners(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	if planners, err := user.QueryPlanners().All(c); err != nil {
		return err
	} else {
		return c.JSON(planners)
	}
}

func (h *PlannerHandler) CreatePlanner(c fiber.Ctx) error {
	user := c.Locals("user").(*ent.User)

	data := new(models.CreatePlannerReq)
	if err := c.Bind().Body(data); err != nil {
		return fiber.ErrBadRequest
	}

	planner, err := h.db.Planner.Create().
		SetType(data.Type).
		SetName(data.Name).
		SetUser(user).
		Save(c)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(planner)
}

func (h *PlannerHandler) UpdatePlanner(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	data := new(models.RenamePlannerReq)

	if err := c.Bind().Body(data); err != nil {
		return fiber.ErrBadRequest
	}

	planner.Update().SetName(data.Name).Exec(c)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PlannerHandler) DeletePlanner(c fiber.Ctx) error {
	planner := c.Locals("planner").(*ent.Planner)

	if err := h.db.Planner.DeleteOne(planner).Exec(c); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
