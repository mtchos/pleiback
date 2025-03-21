package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mtchos/pleiback/internal/adapter/handler"
	"github.com/mtchos/pleiback/internal/adapter/integration/spotify"
	"github.com/mtchos/pleiback/internal/config"
	"github.com/mtchos/pleiback/internal/domain/service"
	"log/slog"
	"net/http"
)

func main() {
	config.SetupLogger()

	r := chi.NewRouter()

	r.Use(config.ContextMiddleware)
	r.Use(config.LoggingMiddleware)

	httpClient := &http.Client{}
	spotifyService := spotify.NewClient(httpClient)
	trackService := service.NewTrackService(spotifyService)
	trackHandler := handler.NewTrackHandler(trackService)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello mundo"))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})

	r.Mount("/tracks", trackHandler.Routes())

	port := "8080"
	slog.Info("server running on", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("error", err)
		return
	}
}
