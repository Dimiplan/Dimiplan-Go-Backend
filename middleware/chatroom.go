package middleware

import (
	"dimiplan-backend/ent"
	"dimiplan-backend/ent/chatroom"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func QueryChatroomMiddleware(db *ent.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		owner := c.Locals("user").(*ent.User)
		chatroomID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return fiber.ErrBadRequest
		}
		chatroom, err := owner.QueryOwnedChatrooms().Where(chatroom.ID(chatroomID)).Only(c)
		if err != nil {
			if chatroom == nil {
				return fiber.ErrBadRequest
			}
			return err
		}
		c.Locals("chatroom", chatroom)
		return c.Next()
	}
}
