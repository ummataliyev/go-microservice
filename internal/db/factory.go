package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/thealish/go-microservice/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// NewDatabase returns a GORM handle for the SQL database specified by
// cfg.Server.DBProvider. Supported providers: "postgres", "mysql".
func NewDatabase(ctx context.Context, cfg config.Config, log zerolog.Logger) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Server.DBProvider) {
	case "postgres":
		return NewPostgres(ctx, cfg.Postgres, log)
	case "mysql":
		return NewMySQL(ctx, cfg.MySQL, log)
	default:
		return nil, fmt.Errorf("unsupported SQL db_provider: %s", cfg.Server.DBProvider)
	}
}

// NewMongoDatabase returns a MongoDB database handle using the Mongo config.
func NewMongoDatabase(ctx context.Context, cfg config.Config, log zerolog.Logger) (*mongo.Database, error) {
	return NewMongo(ctx, cfg.Mongo, log)
}
