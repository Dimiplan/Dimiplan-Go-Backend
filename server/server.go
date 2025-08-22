package server

import (
	"dimiplan-backend/config"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
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

	app.Use(helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "cross-origin",
		CrossOriginEmbedderPolicy: "credentialless",
		ContentSecurityPolicy: `default-src 'self';
		img-src 'self' data: https://*.googleusercontent.com;
		style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com;
		script-src 'self' 'unsafe-inline' https://static.cloudflareinsights.com https://*.cloudflare.com;
		font-src 'self' https://fonts.gstatic.com https://cdn.jsdelivr.net`,
	}))

	//	accessLog, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//	if err != nil {
	//		panic(err)
	//	}

	app.Use(logger.New(logger.Config{
		// Format:     "{time: \"${time}\", ip: \"${ip}\", method: \"${method}\", url: \"${url}\", status: \"${status}\", error: \"${error}\"}\n",
		TimeFormat: "01-02 15:04:05",
		TimeZone:   "Asia/Seoul",
		//		Stream:     accessLog,
	}))
	app.Use(compress.New())

	storage := redis.New(*cfg.RedisConfig)

	app.Use(session.New(session.Config{
		Storage:        storage,
		Extractor:      session.FromCookie("dimiplan.sid"),
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "lax",
		IdleTimeout:    time.Hour * 24,
	}))

	app.Use("/", static.New("/", static.Config{
		FS:            os.DirFS("dist"),
		Compress:      true,
		CacheDuration: time.Hour * 24,
		MaxAge:        60 * 60 * 12,
	}))

	app.Use(func(c fiber.Ctx) error {
		if c.Protocol() == "http" {
			return c.Redirect().To(fmt.Sprintf("https://%s%s", c.Hostname(), c.Path()))
		}
		return c.Next()
	})

	return app
}
