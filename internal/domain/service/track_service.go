package service

import (
	"github.com/mtchos/pleiback/internal/adapter/integration/spotify"
	"github.com/mtchos/pleiback/internal/domain/entity"
	"github.com/mtchos/pleiback/internal/domain/usecase"
	"log/slog"
)

type trackServiceImpl struct {
	spotify spotify.Integration
}

func NewTrackService(spotify spotify.Integration) usecase.TrackService {
	return &trackServiceImpl{spotify: spotify}
}

func (s *trackServiceImpl) Search(query string) ([]entity.Track, error) {
	limit := int64(10)
	offset := int64(0)

	tracks, err := s.spotify.SearchTracks(query, limit, offset)
	if err != nil {
		slog.Error("usecase error when searching track", "error", err)
		return nil, err
	}

	return tracks, nil
}
