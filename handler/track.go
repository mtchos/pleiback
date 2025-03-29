package handler

import (
	"context"
	"encoding/json"
	"github.com/mtchos/pleiback/service"
	"github.com/mtchos/pleiback/spotify"
	"log/slog"
	"net/http"
	"strconv"
)

type Spotify interface {
	SearchTracks(ctx context.Context, filter spotify.SearchTrackFilter) (*spotify.TracksResponse, error)
}

type Track struct {
	spotify Spotify
}

func TrackInstance(spotify Spotify) *Track {
	return &Track{
		spotify: spotify,
	}
}

func (h *Track) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	response, err := h.spotify.SearchTracks(ctx, spotify.SearchTrackFilter{
		Query:  query,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		slog.Error("error retrieving spotify response", "err", err)
		http.Error(w, "error retrieving response from spotify", http.StatusInternalServerError)
		return
	}
	if response.Tracks.Items == nil {
		slog.Error("error retrieving spotify response", "err", err)
		http.Error(w, "no tracks retrieved from spotify", http.StatusNotFound)
		return
	}

	var tracks []*service.Track
	for _, t := range response.Tracks.Items {

		var artists []*service.Artist
		for _, a := range t.Artists {
			artists = append(artists, &service.Artist{
				ID:     a.ID,
				Name:   a.Name,
				Href:   a.Href,
				URI:    a.URI,
				Genres: []*string{},
			})
		}

		tracks = append(tracks, &service.Track{
			ID:      t.ID,
			Name:    t.Name,
			Href:    t.Href,
			URI:     t.URI,
			Artists: artists,
		})
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&tracks); err != nil {
		slog.Error("error writing response message", "err", err)
		http.Error(w, "error writing response message", http.StatusInternalServerError)
		return
	}
}
