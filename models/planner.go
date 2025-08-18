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

type CreateTaskReq struct {
	ID       int    `json:"planner_id"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type UpdateTaskReq struct {
	ID       int    `uri:"id"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type DeleteTaskReq struct {
	ID int `uri:"id"`
}
