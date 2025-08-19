package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func QueryChatroomMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		user := c.Locals("user").(*ent.User)
		chatroomID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid chatroom ID",
			})
		}
		chatroom, err := user.QueryChatrooms().Where(chatroom.ID(chatroomID)).Only(c)
		if err != nil {
			log.Error(err)
			if chatroom == nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Chatroom not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve chatroom",
			})
		}
		c.Locals("chatroom", chatroom)
		return c.Next()
	}
}
