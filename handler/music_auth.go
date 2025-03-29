package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type SpotifyUserAuth interface {
	SetToken(w http.ResponseWriter, token *string, expiresAt *time.Time) error
	ExchangeCodeForToken(code string) (*string, *time.Time, error)
}

type MusicAuth struct {
	spotifyUserAuth SpotifyUserAuth
	url             string
}

func MusicAuthInstance(spotifyUserAuth SpotifyUserAuth) *MusicAuth {
	return &MusicAuth{
		spotifyUserAuth: spotifyUserAuth,
		url:             "https://accounts.spotify.com",
	}
}

func (h *MusicAuth) Redirect(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	scopes := "playlist-modify-public"

	authURL := fmt.Sprintf("%s/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		h.url, clientID, redirectURI, scopes)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *MusicAuth) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "authorization code not found", http.StatusInternalServerError)
	}

	token, expiresAt, err := h.spotifyUserAuth.ExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.spotifyUserAuth.SetToken(w, token, expiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
