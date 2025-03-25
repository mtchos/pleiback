package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/mtchos/pleiback/old/internal/domain/usecase"
	"log/slog"
	"net/http"
)

type TrackHandler struct {
	service usecase.TrackService
}

func NewTrackHandler(service usecase.TrackService) *TrackHandler {
	return &TrackHandler{service: service}
}

func (h *TrackHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.Search)
	return r
}

func (h *TrackHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		slog.Error("query parameter 'q' is required")
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	tracks, err := h.service.Search(query)
	if err != nil {
		slog.Error("handler error when searching tracks", "error", err)
		http.Error(w, "error searching tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tracks)
	if err != nil {
		slog.Error("error encoding search tracks response", "error", err)
		http.Error(w, "error enconding tracks", http.StatusInternalServerError)
	}
}
