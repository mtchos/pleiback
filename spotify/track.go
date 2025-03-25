package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mtchos/pleiback/track"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type Track struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PreviewURL string `json:"preview_url"`
	URI        string `json:"uri"`

	Album struct {
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name string `json:"name"`
	} `json:"album"`

	Artists []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"artists"`
}

type findTracksResponse struct {
	Tracks struct {
		Total int      `json:"total"`
		Items []*Track `json:"items"`
	} `json:"tracks"`
}

func (c *Client) FindTracks(ctx context.Context, query string, limit, offset int) ([]*track.Track, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.URL+"/search", nil)
	if err != nil {
		slog.Error("could not retrieve tracks from spotify api", "err", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	err = godotenv.Load(".env")
	if err != nil {
		slog.Error("could not load .env file", "err", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("could not retrieve tracks from spotify api", "err", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.Error("spotify api response status code is not 200", "status", resp.StatusCode)
		return nil, fmt.Errorf("response status code is %s", resp.Status)
	}

	var tracksResponse findTracksResponse
	if err = json.NewDecoder(resp.Body).Decode(&tracksResponse); err != nil {
		slog.Error("could not decode response body", "err", err)
		return nil, err
	}

	var tracks []*track.Track
	for _, spotifyTrack := range tracksResponse.Tracks.Items {
		tracks = append(tracks, &track.Track{
			ID:     spotifyTrack.ID,
			Name:   spotifyTrack.Name,
			Artist: spotifyTrack.Artists[0].Name,
		})
	}

	if err = resp.Body.Close(); err != nil {
		slog.Error("could not close response body", "err", err)
		return nil, err
	}

	return tracks, nil
}
