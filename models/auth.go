package models

type GoogleResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}
