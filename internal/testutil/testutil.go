package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go-microservice/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

	sqlDB, err := db.DB()
	if err != nil {
		t.Skipf("skipping integration test: unable to get underlying sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("skipping integration test: postgres ping failed: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto-migrate failed: %v", err)
	}

	return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		t.Errorf("failed to truncate users table: %v", err)
	}
}

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
