package auth

import (
	"dimiplan-backend/models"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3/client"
)

func GetUser(token string) models.GoogleResponse {
	cc := client.New()
	res, err := cc.Get("https://www.googleapis.com/oauth2/v1/userinfo", client.Config{
		Header: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		},
	})
	if err != nil {
		panic(err)
	}
	var data models.GoogleResponse
	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		panic(err)
	}
	return data
}
