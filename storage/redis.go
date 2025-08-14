package storage

import (
	"context"
	"dimiplan-backend/models"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

func (r *RedisService) SaveUser(user *models.User) error {
	ctx := context.Background()
	userJSON, _ := json.Marshal(user)
	return r.client.Set(ctx, "user:"+user.ID, userJSON, 24*time.Hour).Err()
}

func (r *RedisService) GetUser(userID string) (*models.User, error) {
	ctx := context.Background()
	userJSON, err := r.client.Get(ctx, "user:"+userID).Result()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	return &user, err
}

func (r *RedisService) DeleteUser(userID string) error {
	ctx := context.Background()
	return r.client.Del(ctx, "user:"+userID).Err()
}
