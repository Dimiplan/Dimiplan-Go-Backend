package handlers

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"

	"golang.org/x/oauth2"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	oauth      *oauth2.Config
	sessionSvc *auth.SessionService
	db         *ent.Client
}

func NewAuthHandler(oauth *oauth2.Config, sessionSvc *auth.SessionService, db *ent.Client) *AuthHandler {
	return &AuthHandler{
		oauth:      oauth,
		sessionSvc: sessionSvc,
		db:         db,
	}
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.oauth.AuthCodeURL("state")
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	token, error := h.oauth.Exchange(c.Context(), c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	data := auth.GetUser(token.AccessToken)
	h.sessionSvc.SetIDInSession(c, data.ID)
	user, err := h.db.User.Query().Where(user.ID(data.ID)).First(c.Context())
	if user == nil {
		h.db.User.Create().SetID(data.ID).SetEmail(data.Email).SetName(data.Name).SetProfileURL(data.ProfileURL).SaveX(c.Context())
	} else if err != nil {
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged in successfully"})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.sessionSvc.ClearSession(c)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
