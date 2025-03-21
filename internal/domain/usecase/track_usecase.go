package usecase

import (
	"github.com/mtchos/pleiback/internal/domain/entity"
)

type TrackService interface {
	Search(query string) ([]entity.Track, error)
}
