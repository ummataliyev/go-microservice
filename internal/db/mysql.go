package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/thealish/go-microservice/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQL opens a GORM connection to MySQL, configures the pool, and pings
// the database to verify connectivity.
func NewMySQL(ctx context.Context, cfg config.MySQLConfig, log zerolog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.PoolSize)
	sqlDB.SetMaxIdleConns(cfg.PoolMaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.PoolMaxLifetime)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	log.Info().Msg("connected to MySQL")
	return db, nil
}
