package models

type CreatePlannerReq struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type GetTasksReq struct {
	ID int `uri:"id"`
}

type RenamePlannerReq struct {
	ID   int    `uri:"id"`
	Name string `json:"name"`
}

type DeletePlannerReq struct {
	ID int `uri:"id"`
}
