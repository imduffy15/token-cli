package client

import (
	"encoding/json"
	"net/http"
	"strconv"
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
		return Token{}, parseError(err.Error(), bytes)
	}

	return Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		IDToken:      resp.IDToken,
		TokenType:    resp.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).Unix(),
	}, nil
}

type AuthorizationCodeClient struct {
	ClientID     string
	ClientSecret string
}

type RefreshTokenClient struct {
	ClientID     string
	ClientSecret string
}

func (acc AuthorizationCodeClient) RequestToken(httpClient *http.Client, config Config, code string, redirectURI string) (Token, error) {
	body := map[string]string{
		"grant_type":    string(AuthCode),
		"client_id":     acc.ClientID,
		"client_secret": acc.ClientSecret,
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
		"client_secret": rc.ClientSecret,
		"response_type": "token",
	}

	return postToOAuthToken(httpClient, config, body)
}

type GrantType string

const (
	RefreshToken = GrantType("refresh_token")
	AuthCode     = GrantType("authorization_code")
)

type FlexInt int64

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int64)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}

type tokenResponse struct {
	AccessToken  string  `json:"access_token"`
	IDToken      string  `json:"id_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    FlexInt `json:"expires_in"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
}
