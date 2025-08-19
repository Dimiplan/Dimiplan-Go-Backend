package models

type CreateTaskReq struct {
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

type UpdateTaskReq struct {
	TaskID   int    `uri:"task"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}
