package handlers

import (
	"dimiplan-backend/auth"
	"dimiplan-backend/storage"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	oauth   *auth.OAuthService
	jwt     *auth.JWTService
	storage *storage.RedisService
}

func NewAuthHandler(oauth *auth.OAuthService, jwt *auth.JWTService, storage *storage.RedisService) *AuthHandler {
	return &AuthHandler{
		oauth:   oauth,
		jwt:     jwt,
		storage: storage,
	}
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state, err := h.oauth.GenerateState()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate state",
		})
	}

	err = h.oauth.SaveState(state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save state",
		})
	}

	url := h.oauth.GetAuthURL(state)
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "State parameter missing",
		})
	}

	err := h.oauth.ValidateState(state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid state",
		})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code missing",
		})
	}

	token, err := h.oauth.ExchangeCode(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token exchange failed",
		})
	}

	user, err := h.oauth.GetUserInfo(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}

	h.storage.SaveUser(user)

	jwtToken, err := h.jwt.Generate(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate JWT",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   jwtToken,
		"user":    user,
	})
}
