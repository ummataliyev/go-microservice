package db

import (
	"context"

	"go-microservice/internal/config"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewDatabase(ctx context.Context, cfg config.Config, log zerolog.Logger) (*gorm.DB, error) {
	return NewPostgres(ctx, cfg.Postgres, log)
}
