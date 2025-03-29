package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type Auth interface {
	GetToken() (*string, error)
}

type Service struct {
	client *http.Client
	auth   Auth
	url    string
}

func NewService(client *http.Client, authClient Auth) *Service {
	return &Service{
		client: client,
		auth:   authClient,
		url:    os.Getenv("SPOTIFY_API_URL"),
	}
}

func (s *Service) SearchTracks(ctx context.Context, filter SearchTrackFilter) (*TracksResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url+"/search", nil)
	if err != nil {
		slog.Error("error creating request", "err", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", filter.Query)
	q.Add("offset", strconv.Itoa(filter.Offset))
	q.Add("limit", strconv.Itoa(filter.Limit))
	q.Add("type", "track")
	req.URL.RawQuery = q.Encode()

	token, err := s.auth.GetToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+*token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		slog.Error("status response error", "err", err)
		return nil, errors.New("spotify response status code was" + resp.Status)
	}

	var tracks TracksResponse
	if err := json.NewDecoder(resp.Body).Decode(&tracks); err != nil {
		slog.Error("error decoding spotify response", "err", err)
		return nil, err
	}

	return &tracks, nil
}
