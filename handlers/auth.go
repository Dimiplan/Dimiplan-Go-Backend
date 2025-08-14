package handlers

import (
	"dimiplan-backend/auth"

	"golang.org/x/oauth2"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	oauth      *oauth2.Config
	sessionSvc *auth.SessionService
}

func NewAuthHandler(oauth *oauth2.Config, sessionSvc *auth.SessionService) *AuthHandler {
	return &AuthHandler{
		oauth:      oauth,
		sessionSvc: sessionSvc,
	}
}

// Auth fiber handler
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.oauth.AuthCodeURL("state")
	return c.Redirect(url)
}

// Callback to receive google's response
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	token, error := h.oauth.Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	email := auth.GetEmail(token.AccessToken)
	return c.Status(200).JSON(fiber.Map{"email": email, "login": true})
}
