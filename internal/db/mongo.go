package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/the_elita/go-microservice/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongo connects to MongoDB using the provided URI, pings the server, and
// returns a handle to the configured database.
func NewMongo(ctx context.Context, cfg config.MongoConfig, log zerolog.Logger) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	log.Info().Str("database", cfg.Database).Msg("connected to MongoDB")
	return client.Database(cfg.Database), nil
}
