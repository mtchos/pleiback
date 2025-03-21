package spotify

import (
	"github.com/mtchos/pleiback/internal/domain/entity"
)

type Integration interface {
	SearchTracks(query string, limit, offset int64) (SearchTracksResponse, error)
	GetArtists(artistsIDs []string) (GetArtistsResponse, error)
	CreatePlaylist(userID string, playlist entity.Playlist) (CreatePlaylistResponse, error)
	AddPlaylistTracks(playlistID string, tracksIDs []string) error
}

type SearchTracksResponse = []entity.Track

var SearchTracksResult struct {
	Tracks struct {
		Items []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name string `json:"name"`
			} `json:"album"`
			URI string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

type GetArtistsResponse struct {
	Artists []entity.Artist
}

type CreatePlaylistResponse struct {
	Tracks      []entity.Track
	Name        string
	Description string
	URI         string
}
