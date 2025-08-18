package models

type CreatePlannerReq struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type GetTasksInPlannerReq struct {
	Id int `uri:"id"`
}
