package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	key = "spotify_access_token"
)

type Redis interface {
	SetToken(key, token *string, expiresAt *time.Time) error
	GetToken(key string) (string, error)
}

type AuthUserService struct {
	client       *http.Client
	redis        Redis
	url          string
	redirectURI  string
	clientID     string
	clientSecret string
}

func NewAuthUserService(client *http.Client) *AuthUserService {
	return &AuthUserService{
		client:       client,
		url:          os.Getenv("SPOTIFY_AUTH_API_URL"),
		redirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		clientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}

func (s *AuthUserService) SetToken(w http.ResponseWriter, token *string, expiresAt *time.Time) error {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    *token,
		Expires:  *expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return s.redis.SetToken(&key, token, expiresAt)
}

func (s *AuthUserService) GetToken(r *http.Request) (*string, error) {
	cookie, err := r.Cookie(key)
	if err == nil {
		return &cookie.Value, nil
	}

	token, err := s.redis.GetToken(key)
	if err == nil {
		return &token, nil
	}

	return &token, nil
}

func (s *AuthUserService) ExchangeCodeForToken(code string) (*string, *time.Time, error) {
	data := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s",
		code, s.redirectURI)

	req, err := http.NewRequest(http.MethodPost, s.url+"/token", bytes.NewBufferString(data))
	if err != nil {
		slog.Error("error creating request to exchange status for code", "err", err)
		return nil, nil, err
	}

	req.SetBasicAuth(s.clientID, s.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("error requesting token from code", "err", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("status code was not ok", "status", resp.StatusCode)
		return nil, nil, errors.New("failed to exchange status for code")
	}

	var tk TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		slog.Error("error parsing token response", "err", err)
		return nil, nil, err
	}

	expiresAt := time.Now().Add(time.Duration(tk.ExpiresIn))
	return &tk.AccessToken, &expiresAt, nil
}
