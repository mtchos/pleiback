package spotify

import "net/http"

type Client struct {
	client *http.Client
	URL    string
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
		URL:    "https://api.spotify.com/v1",
	}
}
