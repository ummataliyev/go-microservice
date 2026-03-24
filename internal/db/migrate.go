package db

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies all pending SQL migrations from the embedded
// filesystem. It treats migrate.ErrNoChange as a no-op rather than an error.
func RunMigrations(dbURL string, log zerolog.Logger) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("create iofs source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("migrations: no changes to apply")
			return nil
		}
		return fmt.Errorf("run migrations up: %w", err)
	}

	log.Info().Msg("migrations applied successfully")
	return nil
}
