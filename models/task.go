package models

type GetTasksReq struct {
	PlannerID int `uri:"planner"`
}

type CreateTaskReq struct {
	PlannerID int    `uri:"planner"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
}

type UpdateTaskReq struct {
	PlannerID int    `uri:"planner"`
	TaskID    int    `uri:"task"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
}

type DeleteTaskReq struct {
	PlannerID int `uri:"planner"`
	TaskID    int `uri:"task"`
}
