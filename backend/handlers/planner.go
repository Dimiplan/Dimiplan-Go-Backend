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

func (h *PlannerHandler) GetPlanners(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)

	if planners, err := user.QueryPlanners().All(c); err != nil {
		return nil, err
	} else {
		return models.GetPlannersResponse{Planners: planners}, nil
	}
}

func (h *PlannerHandler) CreatePlanner(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	user := c.Locals("user").(*ent.User)

	request := rawRequest.(*models.CreatePlannerRequest)

	_, err := h.db.Planner.Create().
		SetType(request.Type).
		SetName(request.Name).
		SetUser(user).
		Save(c)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *PlannerHandler) UpdatePlanner(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	planner := c.Locals("planner").(*ent.Planner)

	request := rawRequest.(*models.RenamePlannerRequest)

	if err := c.Bind().Body(request); err != nil {
		return nil, fiber.ErrBadRequest
	}

	planner.Update().SetName(request.Name).Exec(c)
	return nil, nil
}

func (h *PlannerHandler) DeletePlanner(rawRequest interface{}, c fiber.Ctx) (interface{}, error) {
	planner := c.Locals("planner").(*ent.Planner)

	if err := h.db.Planner.DeleteOne(planner).Exec(c); err != nil {
		return nil, err
	}
	return nil, nil
}
