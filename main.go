package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/mtchos/pleiback/handler"
	"github.com/mtchos/pleiback/spotify"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("error loading godotenv", "err", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	stfyUserAuth := spotify.NewAuthUserService(client)
	stfyAuth := spotify.NewAuthService(client)
	stfy := spotify.NewService(client, stfyAuth)

	track := handler.TrackInstance(stfy)
	r.Get("/tracks", track.Find)

	musicAuth := handler.MusicAuthInstance(stfyUserAuth)
	r.Get("/music/auth/redirect", musicAuth.Redirect)
	r.Get("/music/auth/callback", musicAuth.Callback)

	slog.Info("server runnning on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("server has stopped running")
		return
	}
}
