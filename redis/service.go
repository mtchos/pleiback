package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log/slog"
	"os"
	"time"
)

type Service struct {
	client *redis.Client
}

func NewService() *Service {
	return &Service{
		client: redis.NewClient(
			&redis.Options{
				Addr:     os.Getenv("REDIS_ADDR"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       0,
			},
		),
	}
}

func (s *Service) SetToken(key string, value *string, expiresAt *time.Time) error {
	return s.client.Set(context.Background(), key, *value, time.Until(*expiresAt)).Err()
}

func (s *Service) GetToken(value string) (*string, error) {
	token, err := s.client.Get(context.Background(), value).Result()
	if errors.Is(err, redis.Nil) {
		slog.Error("no token found in redis")
		return nil, errors.New("no token found in redis")
	} else if err != nil {
		slog.Error("error retrieving redis token")
		return nil, err
	}

	return &token, nil
}
