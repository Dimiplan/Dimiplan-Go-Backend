package models

import "dimiplan-backend/ent"

type GetPlannersResponse struct {
	Planners []*ent.Planner `json:"planners"`
}

type CreatePlannerRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RenamePlannerRequest struct {
	PlannerID int    `path:"planner"`
	Name      string `json:"name"`
}

type DeletePlannerRequest struct {
	PlannerID int `path:"planner"`
}
