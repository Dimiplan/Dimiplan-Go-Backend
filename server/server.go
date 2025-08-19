package server

import (
	"dimiplan-backend/config"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/storage/redis/v3"
)

func Setup(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(helmet.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	app.Use(compress.New())

	storage := redis.New(*cfg.RedisConfig)

	app.Use(session.New(session.Config{
		Storage:   storage,
		Extractor: session.FromCookie("dimiplan.sid"),
	}))

	app.Use("/", static.New("/", static.Config{
		FS:       os.DirFS("dist"),
		Compress: true,
	}))

	return app
}
