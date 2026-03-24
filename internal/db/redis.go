package db

import (
	"context"
	"fmt"

	"go-microservice/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedis(ctx context.Context, cfg config.RedisConfig, log zerolog.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:    cfg.Password,
		DB:          cfg.DB,
		DialTimeout: cfg.Timeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Warn().Err(err).Msg("redis ping failed — running without cache")
		return nil
	}

	log.Info().Msg("connected to Redis")
	return client
}
