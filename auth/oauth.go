package auth

import (
	"context"
	"crypto/rand"
	"dimiplan-backend/models"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config *oauth2.Config
	redis  *redis.Client
}

func NewOAuthService(config *oauth2.Config, redisClient *redis.Client) *OAuthService {
	return &OAuthService{
		config: config,
		redis:  redisClient,
	}
}

func (o *OAuthService) GenerateState() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (o *OAuthService) SaveState(state string) error {
	ctx := context.Background()
	stateData := models.OAuthState{
		State:     state,
		CreatedAt: time.Now(),
	}
	stateJSON, _ := json.Marshal(stateData)
	return o.redis.Set(ctx, "oauth_state:"+state, stateJSON, 5*time.Minute).Err()
}

func (o *OAuthService) ValidateState(state string) error {
	ctx := context.Background()
	_, err := o.redis.Get(ctx, "oauth_state:"+state).Result()
	if err != nil {
		return err
	}
	o.redis.Del(ctx, "oauth_state:"+state)
	return nil
}

func (o *OAuthService) GetAuthURL(state string) string {
	return o.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (o *OAuthService) ExchangeCode(code string) (*oauth2.Token, error) {
	return o.config.Exchange(context.Background(), code)
}

func (o *OAuthService) GetUserInfo(token *oauth2.Token) (*models.User, error) {
	ctx := context.Background()
	client := o.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
