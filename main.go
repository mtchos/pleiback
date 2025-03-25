package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mtchos/pleiback/spotify"
	"github.com/mtchos/pleiback/track"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	client := &http.Client{}

	spotifyClient := spotify.NewClient(client)

	tracks := track.NewTracks(spotifyClient)
	trackHandler := track.NewHandler(tracks)
	trackHandler.Routes(r)

	port := "8080"
	slog.Info("server running on", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
