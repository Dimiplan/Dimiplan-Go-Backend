package handlers

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"

	"golang.org/x/oauth2"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type AuthHandler struct {
	oauth *oauth2.Config
	db    *ent.Client
}

func NewAuthHandler(oauth *oauth2.Config, db *ent.Client) *AuthHandler {
	return &AuthHandler{
		oauth: oauth,
		db:    db,
	}
}

func (h *AuthHandler) GoogleLogin(c fiber.Ctx) error {
	url := h.oauth.AuthCodeURL("state")
	return c.Redirect().To(url)
}

func (h *AuthHandler) GoogleCallback(c fiber.Ctx) error {
	token, error := h.oauth.Exchange(c, c.FormValue("code"))
	if error != nil {
		panic(error)
	}
	data := auth.GetUser(token.AccessToken)
	session.FromContext(c).Set("id", data.ID)
	user, err := h.db.User.Query().Where(user.ID(data.ID)).First(c)
	if user == nil {
		h.db.User.Create().SetID(data.ID).SetEmail(data.Email).SetName(data.Name).SetProfileURL(data.ProfileURL).SaveX(c)
	} else if err != nil {
		panic(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged in successfully"})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	session.FromContext(c).Destroy()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
}
