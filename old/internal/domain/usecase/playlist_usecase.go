package usecase

import (
	entity2 "github.com/mtchos/pleiback/old/internal/domain/entity"
)

type PlaylistService interface {
	GeneratePlaylist(track entity2.Track) (entity2.Playlist, error)
	CreatePlaylist(userID string, tracks []entity2.Track) (entity2.Playlist, error)
}
