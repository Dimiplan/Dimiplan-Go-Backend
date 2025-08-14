package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:8080",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(compress.New())

	app.Static("/", "dist")

	return app
}
