package auth

import (
	"encoding/json"
	"errors"
	"key-vault-server/config"
	"net/http"
	"net/url"
	"strings"
)

const (
	accessTokenType   = "access_token"
	introspectionPath = "/oauth2/introspect"
)

var (
	wrongTokenType = errors.New("wrong token type")
	inactiveToken  = errors.New("inactive token")
)

type introspectionResponse struct {
	Active    bool   `json:"active"`
	TokenType string `json:"token_type"`
	UserId    string `json:"sub"`
}

func getUserId(cfg *config.Config, token *string) (*string, error) {
	// Prepare query
	headers := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"application/json"},
	}

	body := url.Values{}
	body.Set("token", *token)

	// Run query
	req, err := http.NewRequest("POST", cfg.HydraAdminUrl+introspectionPath, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = headers

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	// Decode response
	var decodedResponse introspectionResponse
	err = json.NewDecoder(response.Body).Decode(&decodedResponse)
	if err != nil {
		return nil, err
	}

	// Evaluate token information
	if decodedResponse.TokenType != accessTokenType {
		return nil, wrongTokenType
	}
	if !decodedResponse.Active && len(decodedResponse.UserId) != 0 {
		return nil, inactiveToken
	}
	return &decodedResponse.UserId, nil
}
