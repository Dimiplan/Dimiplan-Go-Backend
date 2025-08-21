package models

import "dimiplan-backend/ent"

type GetTasksRequest struct {
	PlannerID int `path:"planner"`
}

type GetTasksResponse struct {
	Tasks []*ent.Task `json:"tasks"`
}

type CreateTaskRequest struct {
	PlannerID int    `path:"planner"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
}

type UpdateTaskRequest struct {
	PlannerID int    `path:"planner"`
	TaskID    int    `path:"task"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
}

type DeleteTaskRequest struct {
	PlannerID int `path:"planner"`
	TaskID    int `path:"task"`
}
