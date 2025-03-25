package spotify

import (
	"net/http"
)

type AuthClient struct {
	client  *http.Client
	baseURL string
}

func NewAuthClient(client *http.Client) AuthIntegration {
	return &AuthClient{
		client:  client,
		baseURL: "https://accounts.spotify.com",
	}
}

func (s AuthClient) Authorize() (AuthorizeResponse, error) {
	//ctx := context.Background()
	//req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL, nil)
	return AuthorizeResponse{}, nil
}

func (s AuthClient) GetToken(code string) (GetTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
