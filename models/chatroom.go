package models

type ChatroomID struct {
	ID int `uri:"id"`
}

type UpdateChatroom struct {
	ID int `uri:"id"`
	Name string `json:"name"`
}
