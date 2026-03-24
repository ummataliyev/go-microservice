package db

import (
	"context"
	"fmt"

	"go-microservice/internal/config"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(ctx context.Context, cfg config.PostgresConfig, log zerolog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.PoolSize)
	sqlDB.SetMaxIdleConns(cfg.PoolMaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.PoolMaxLifetime)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	log.Info().Msg("connected to PostgreSQL")
	return db, nil
}
