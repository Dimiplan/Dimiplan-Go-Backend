package auth

import (
	"dimiplan-backend/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetEmail of user
func GetUser(token string) models.GoogleResponse {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		panic(err)
	}
	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken},
		},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data models.GoogleResponse
	fmt.Println(string(body))
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return data
}
