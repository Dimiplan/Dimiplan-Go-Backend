package models

import (
	"time"

	"dimiplan-backend/ent"
)

// --- DTOs ---
// UserRes is a response DTO to prevent over-posting and to control output shape.
type UserRes struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	ProfileURL string    `json:"profileURL"`
	Admin      bool      `json:"admin"`
	Plan       string    `json:"plan"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UpdateUserReq struct {
	Name *string `json:"name"`
}

// toUserRes maps ent.User -> UserRes.
func toUserRes(u *ent.User) UserRes {
	return UserRes{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		ProfileURL: u.ProfileURL,
		Admin:      u.Admin,
		Plan:       u.Plan,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
