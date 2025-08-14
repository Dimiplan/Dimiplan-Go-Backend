package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"verified_email"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type OAuthState struct {
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
