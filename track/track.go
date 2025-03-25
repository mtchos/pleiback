package track

import (
	"context"
	"fmt"
	"log/slog"
)

type Track struct {
	ID     string
	Name   string
	Artist string
}

type Spotify interface {
	FindTracks(ctx context.Context, query string, limit, offset int) ([]*Track, error)
}

type Tracks struct {
	spotify Spotify
}

func NewTracks(spotify Spotify) *Tracks {
	return &Tracks{spotify}
}

func (t *Tracks) Find(ctx context.Context, query string) ([]*Track, error) {
	resp, err := t.spotify.FindTracks(ctx, query, 20, 0)
	if err != nil {
		slog.Error("business error when searching tracks", "error", err)
		return nil, err
	}

	if resp == nil {
		slog.Error("no response body when searching tracks", "error", err)
		return nil, fmt.Errorf("no response body")
	}

	return resp, nil
}
