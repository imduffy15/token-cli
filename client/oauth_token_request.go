package client

import (
	"encoding/json"
	"net/http"
	"time"
)

func postToOAuthToken(httpClient *http.Client, config Config, body map[string]string) (Token, error) {
	bytes, err := UnauthenticatedRequester{}.PostForm(httpClient, config, config.GetActiveTarget().TokenEndpoint, body)
	if err != nil {
		return Token{}, err
	}

	resp := tokenResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return Token{}, parseError(config.GetActiveTarget().TokenEndpoint, bytes)
	}

	return Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		TokenType:    resp.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).Unix(),
	}, nil
}

type AuthorizationCodeClient struct {
	ClientID string
}

type RefreshTokenClient struct {
	ClientID string
}

func (acc AuthorizationCodeClient) RequestToken(httpClient *http.Client, config Config, code string, redirectURI string) (Token, error) {
	body := map[string]string{
		"grant_type":    string(AuthCode),
		"client_id":     acc.ClientID,
		"response_type": "token",
		"redirect_uri":  redirectURI,
		"code":          code,
	}

	return postToOAuthToken(httpClient, config, body)
}

func (rc RefreshTokenClient) RequestToken(httpClient *http.Client, config Config, refreshToken string) (Token, error) {
	body := map[string]string{
		"grant_type":    string(RefreshToken),
		"refresh_token": refreshToken,
		"client_id":     rc.ClientID,
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type GrantType string

const (
	RefreshToken = GrantType("refresh_token")
	AuthCode     = GrantType("authorization_code")
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
}
