package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/the_elita/go-microservice/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB connects to a test postgres database using the POSTGRES_DSN env var
// (default: host=localhost port=5432 user=test password=test dbname=test_db sslmode=disable),
// auto-migrates the User model, and returns the *gorm.DB handle.
// The test is skipped if the connection fails.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=test password=test dbname=test_db sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("skipping integration test: unable to connect to postgres: %v", err)
	}

	// Verify the connection is alive.
	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("skipping integration test: unable to get underlying sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("skipping integration test: postgres ping failed: %v", err)
	}

	// Auto-migrate the User model.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto-migrate failed: %v", err)
	}

	return db
}

// CleanupTestDB truncates the users table so each test starts with a clean slate.
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Errorf("failed to truncate users table: %v", err)
	}
}

// SetupTestRedis connects to a test redis instance using REDIS_HOST / REDIS_PORT env vars
// (defaults: localhost / 6379). The test is skipped if the connection fails.
func SetupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("skipping integration test: unable to connect to redis: %v", err)
	}

	return client
}
