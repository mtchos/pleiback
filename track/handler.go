package track

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	tracks Tracks
}

func NewHandler(tracks *Tracks) *Handler {
	return &Handler{*tracks}
}

func (h *Handler) Routes(r *chi.Mux) {
	r.Get("/tracks", h.FindTracks)
}

func (h *Handler) FindTracks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	trackList, err := h.tracks.Find(ctx, query)
	if err != nil {
		http.Error(w, "error searching tracks, "+fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(trackList); err != nil {
		http.Error(w, "error encoding tracks", http.StatusInternalServerError)
		return
	}
}
