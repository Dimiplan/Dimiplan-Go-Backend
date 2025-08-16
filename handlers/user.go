package handlers

import (
	"dimiplan-backend/ent"
	// "github.com/gofiber/fiber/v3/middleware/session"
)

type UserHandler struct {
	db *ent.Client
}

func NewUserHandler(db *ent.Client) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

/// @brief ID를 통해 유저 정보 가져오기
///
// func (h *UserHandler) GetUser(c fiber.Ctx) error {
// 	uid := c.Locals("id").(string)

// 	u, err := h.db.User.Query().Where(user.ID(uid)).Only(c)
// 	if err != nil {
// 		log.Errorf("Failed to retrieve user: %v", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Internal Server Error",
// 		})
// 	}

// 	return c.JSON(toUserRes(u))
// }
