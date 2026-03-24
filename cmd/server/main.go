package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-microservice/internal/api"
	"go-microservice/internal/api/handlers"
	"go-microservice/internal/api/middleware"
	"go-microservice/internal/config"
	"go-microservice/internal/db"
	"go-microservice/internal/logger"
	"go-microservice/internal/repository"
	"go-microservice/internal/security"
	"go-microservice/internal/service"

	_ "go-microservice/docs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	log := logger.New(cfg.Logging, cfg.Server.Environment)

	ctx := context.Background()

	gormDB, err := db.NewPostgres(ctx, cfg.Postgres, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := db.RunMigrations(cfg.Postgres.DSN, log); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	redisClient := db.NewRedis(ctx, cfg.Redis, log)

	repo := repository.NewGORMUser(gormDB)

	hasher := security.NewBcryptHasher()
	jwtSvc := security.NewJWTService(cfg.JWT)
	authSvc := service.NewAuth(repo, jwtSvc, hasher, redisClient, cfg.Auth)
	userSvc := service.NewUsers(repo, hasher)

	authHandler := handlers.NewAuth(authSvc)
	userHandler := handlers.NewUsers(userSvc)
	healthHandler := handlers.NewHealth(cfg.Server.AppName, cfg.Server.AppVersion)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	api.SetupRouter(app, healthHandler, authHandler, userHandler, jwtSvc, *cfg, redisClient, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	go func() {
		log.Info().Str("addr", addr).Msg("starting server")
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("server listen error")
		}
	}()

	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("received shutdown signal")

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}

	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("error closing redis connection")
		}
	}

	sqlDB, err := gormDB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Msg("error closing database connection")
		}
	}

	log.Info().Msg("shutdown complete")
}
