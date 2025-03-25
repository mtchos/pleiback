package usecase

import (
	"github.com/mtchos/pleiback/old/internal/domain/entity"
)

type TrackService interface {
	Search(query string) ([]entity.Track, error)
}
