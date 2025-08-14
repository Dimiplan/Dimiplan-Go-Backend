package handlers

import (
	"dimiplan-backend/auth"
)

type UserHandler struct {
	sessionSvc *auth.SessionService
}

func NewUserHandler(sessionSvc *auth.SessionService) *UserHandler {
	return &UserHandler{
		sessionSvc: sessionSvc,
	}
}

/*
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := h.sessionSvc.GetIDFromSession(c)

	user, err := h.storage.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	h.storage.DeleteUser(userID)

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *UserHandler) Protected(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	return c.JSON(fiber.Map{
		"message": "Access to protected resource",
		"email":   email,
	})
}
*/
