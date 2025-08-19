package models

type CreatePlannerReq struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RenamePlannerReq struct {
	Name string `json:"name"`
}
