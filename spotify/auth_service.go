package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type AuthService struct {
	client       *http.Client
	url          string
	token        *string
	clientID     string
	clientSecret string
	expiresAt    *time.Time
}

func NewAuthService(client *http.Client) *AuthService {
	return &AuthService{
		client:       client,
		url:          os.Getenv("SPOTIFY_AUTH_API_URL"),
		clientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}

func (s *AuthService) GetToken() (*string, error) {
	now := time.Now()
	if s.token != nil && s.expiresAt != nil && s.expiresAt.After(now) {
		return s.token, nil
	}

	accessToken, expiresAt, err := s.requestNewToken()
	if err != nil {
		slog.Error("could not request new access token",
			"err", err)
		return nil, err
	}

	if accessToken == nil {
		slog.Error("spotify access token is invalid",
			"err", err)
		return nil, err
	}

	if expiresAt == nil {
		slog.Error("spotify expires_in is invalid",
			"err", err)
		return nil, err
	}

	s.token = accessToken
	s.expiresAt = expiresAt

	return accessToken, nil
}

func (s *AuthService) requestNewToken() (*string, *time.Time, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)

	req, err := http.NewRequest(http.MethodPost, s.url+"/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		slog.Error("error requesting new client credential", "err", err)
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("could not request a new client credential", "err", err)
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("error requesting new status code", "status", resp.Status)
		return nil, nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("could not decode response body", "err", err)
		return nil, nil, err
	}

	expiresAt := time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)

	return &response.AccessToken, &expiresAt, nil
}
