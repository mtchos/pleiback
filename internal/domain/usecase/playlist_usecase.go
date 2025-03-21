package usecase

import (
	"github.com/mtchos/pleiback/internal/domain/entity"
)

type PlaylistService interface {
	GeneratePlaylist(track entity.Track) (entity.Playlist, error)
	CreatePlaylist(userID string, tracks []entity.Track) (entity.Playlist, error)
}
