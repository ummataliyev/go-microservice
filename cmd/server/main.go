package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go-microservice/internal/api"
	"go-microservice/internal/api/handlers"
	"go-microservice/internal/api/middleware"
	"go-microservice/internal/config"
	"go-microservice/internal/db"
	"go-microservice/internal/logger"
	"go-microservice/internal/repository"
	"go-microservice/internal/security"
	"go-microservice/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func main() {
	// 1. Load configuration.
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger.
	log := logger.New(cfg.Logging, cfg.Server.Environment)

	ctx := context.Background()
	provider := strings.ToLower(cfg.Server.DBProvider)

	// 3. Connect to the database.
	var gormDB *gorm.DB
	var mongoDB *mongo.Database

	if provider == "mongo" {
		mongoDB, err = db.NewMongoDatabase(ctx, *cfg, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to MongoDB")
		}
	} else {
		gormDB, err = db.NewDatabase(ctx, *cfg, log)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}
	}

	// 4. Run migrations (postgres only).
	if provider == "postgres" {
		if err := db.RunMigrations(cfg.Postgres.DSN, log); err != nil {
			log.Fatal().Err(err).Msg("failed to run migrations")
		}
	}

	// 5. Connect Redis (graceful — returns nil on failure).
	redisClient := db.NewRedis(ctx, cfg.Redis, log)

	// 6. Wire dependencies.
	var repo repository.UserRepository
	switch provider {
	case "postgres", "mysql":
		repo = repository.NewGORMUser(gormDB)
	case "mongo":
		repo = repository.NewUserMongo(mongoDB)
	}

	hasher := security.NewBcryptHasher()
	jwtSvc := security.NewJWTService(cfg.JWT)
	authSvc := service.NewAuth(repo, jwtSvc, hasher, redisClient, cfg.Auth)
	userSvc := service.NewUsers(repo, hasher)

	authHandler := handlers.NewAuth(authSvc)
	userHandler := handlers.NewUsers(userSvc)
	healthHandler := handlers.NewHealth(cfg.Server.AppName, cfg.Server.AppVersion)

	// 7. Create Fiber app with custom error handler.
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// 8. Setup router.
	api.SetupRouter(app, healthHandler, authHandler, userHandler, jwtSvc, *cfg, redisClient, log)

	// 9. Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// 10. Start server in a goroutine.
	go func() {
		log.Info().Str("addr", addr).Msg("starting server")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("server listen error")
		}
	}()

	// Block until shutdown signal.
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("received shutdown signal")

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}

	// Close Redis if connected.
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("error closing redis connection")
		}
	}

	// Close database.
	if gormDB != nil {
		sqlDB, err := gormDB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Error().Err(err).Msg("error closing database connection")
			}
		}
	}
	if mongoDB != nil {
		if err := mongoDB.Client().Disconnect(ctx); err != nil {
			log.Error().Err(err).Msg("error disconnecting from MongoDB")
		}
	}

	log.Info().Msg("shutdown complete")
}
