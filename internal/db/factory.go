package db

import (
	"context"

	"github.com/rs/zerolog"
	"go-microservice/internal/config"
	"gorm.io/gorm"
)

// NewDatabase returns a GORM handle for the PostgreSQL database.
func NewDatabase(ctx context.Context, cfg config.Config, log zerolog.Logger) (*gorm.DB, error) {
	return NewPostgres(ctx, cfg.Postgres, log)
}
