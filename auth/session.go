package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

type SessionService struct {
	store *session.Store
}

func NewSessionService(redisOptions *redis.Config) *SessionService {
	storage := redis.New(*redisOptions)
	store := session.New(session.Config{
		Storage:   storage,
		KeyLookup: "cookie:dimiplan.sid",
	})
	return &SessionService{store: store}
}

func (s *SessionService) GetIDFromSession(c *fiber.Ctx) string {
	session, err := s.store.Get(c)
	if err != nil {
		panic(err)
	}
	return session.Get("id").(string)
}

func (s *SessionService) SetIDInSession(c *fiber.Ctx, id string) error {
	session, err := s.store.Get(c)
	if err != nil {
		return err
	}
	session.Set("id", id)
	return session.Save()
}
