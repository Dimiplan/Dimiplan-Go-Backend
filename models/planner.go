package models

type CreatePlannerReq struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RenamePlannerReq struct {
	PlannerID int    `uri:"planner"`
	Name      string `json:"name"`
}

type DeletePlannerReq struct {
	PlannerID int `uri:"planner"`
}
