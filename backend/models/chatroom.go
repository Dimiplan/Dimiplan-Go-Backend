package models

import "dimiplan-backend/ent"

type ListChatroomsResponse struct {
	Chatrooms []*ent.Chatroom `json:"chatrooms"`
}

type CreateChatroomRequest struct {
	Name string `json:"name"`
}

type CreateChatroomResponse struct {
	ID int `json:"chatroom_id"`
}

type GetMessagesRequest struct {
	ID int `path:"id"`
}

type GetMessagesResponse struct {
	Messages []*ent.Message `json:"messages"`
}

type UpdateChatroomRequest struct {
	Name string `json:"name"`
	ID   int    `path:"id"`
}

type RemoveChatroomRequest struct {
	ID int `path:"id"`
}
