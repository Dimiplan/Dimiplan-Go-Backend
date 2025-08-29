package handlers

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/user"
	"dimiplan-backend/models"

	"golang.org/x/oauth2"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
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

func (h *AuthHandler) Login(c fiber.Ctx) error {
	userId := new(models.LoginRequestWithUid)
	err := c.Bind().Body(userId)
	if err != nil || userId.UID == "" {
		url := h.oauth.AuthCodeURL("state")
		return c.Redirect().To(url)
	} else {
		sess := session.FromContext(c)
		sess.Set("id", userId.UID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func (h *AuthHandler) Callback(c fiber.Ctx) error {
	token, err := h.oauth.Exchange(c, c.FormValue("code"))
	if err != nil {
		log.Error(err)
		return fiber.ErrBadRequest
	}

	err, data := auth.GetUser(token.AccessToken)

	if err != nil {
		log.Error(err)
		return fiber.ErrBadRequest
	}

	sess := session.FromContext(c)
	sess.Set("id", data.ID)

	user, err := h.db.User.Query().Where(user.ID(data.ID)).First(c)
	if user == nil {
		if _, err := h.db.User.Create().SetID(data.ID).SetEmail(data.Email).SetName(data.Name).SetProfileURL(data.ProfileURL).Save(c); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return c.Redirect().To("/")
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	session.FromContext(c).Destroy()
	return c.SendStatus(fiber.StatusNoContent)
}
