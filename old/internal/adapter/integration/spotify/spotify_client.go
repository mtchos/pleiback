package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	entity2 "github.com/mtchos/pleiback/old/internal/domain/entity"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(client *http.Client) Integration {
	return &Client{
		client:  client,
		baseURL: "https://api.spotify.com/v1",
	}
}

func (s Client) SearchTracks(query string, limit int64, offset int64) (SearchTracksResponse, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/search", nil)
	if err != nil {
		slog.Error("client error when creating request", "error", err)
		return SearchTracksResponse{}, nil
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("offset", fmt.Sprintf("%d", offset))
	req.URL.RawQuery = q.Encode()

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	slog.Info("SpotifyClient search request outbound", "request", req)

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("client error when calling SpotifyClient API", "error", err)
		return SearchTracksResponse{}, nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("error closing body", "error", err, "body", Body)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		slog.Error("client error when receiving SpotifyClient API response", "error", err)
		return SearchTracksResponse{},
			fmt.Errorf("SpotifyClient API failed with status %d", resp.StatusCode)
	}

	slog.Info("client response received", "status", resp.StatusCode, "response", resp)

	result := SearchTracksResult
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("client error when decoding JSON result", "error", err)
		return SearchTracksResponse{}, err
	}

	slog.Info("tracks JSON found", "tracks", resp.Body)

	var tracks []entity2.Track
	for _, item := range result.Tracks.Items {
		tracks = append(tracks, entity2.Track{
			ID:      item.ID,
			Name:    item.Name,
			Artists: extractArtists(item.Artists),
			URI:     item.URI})
	}

	slog.Info("tracks found", "tracks", tracks)

	return tracks, nil
}

func (s Client) GetArtists(artistsIDs []string) (GetArtistsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Client) CreatePlaylist(userID string, playlist entity2.Playlist) (CreatePlaylistResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Client) AddPlaylistTracks(playlistID string, tracksIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func extractArtists(artists []struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}) []entity2.Artist {
	var artistEntities []entity2.Artist
	for _, artist := range artists {
		artistEntities = append(artistEntities, entity2.Artist{
			ID:   artist.ID,
			Name: artist.Name,
		})
	}

	return artistEntities
}
