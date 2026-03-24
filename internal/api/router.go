package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/the_elita/go-microservice/internal/api/handlers"
	"github.com/the_elita/go-microservice/internal/api/middleware"
	"github.com/the_elita/go-microservice/internal/config"
	"github.com/the_elita/go-microservice/internal/security"
)

// SetupRouter registers all middleware and routes on the given Fiber app.
func SetupRouter(
	app *fiber.App,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	tokenSvc security.TokenService,
	cfg config.Config,
	redisClient *redis.Client,
	log zerolog.Logger,
) {
	// Global middleware (order matters).
	app.Use(middleware.RequestID())
	app.Use(middleware.Timing(int64(cfg.Logging.SlowRequestThresholdMS)))
	app.Use(middleware.SecurityHeaders())

	// Trusted hosts.
	var trustedHosts []string
	if cfg.Server.TrustedHosts != "" && cfg.Server.TrustedHosts != "*" {
		trustedHosts = strings.Split(cfg.Server.TrustedHosts, ",")
	}
	app.Use(middleware.TrustedHost(trustedHosts))

	app.Use(middleware.CORS(cfg.CORS))

	// Rate limiter.
	rl := middleware.NewRateLimiter(redisClient, cfg.RateLimit)
	app.Use(rl.Middleware())

	// Health / info routes (no auth required).
	app.Get("/", healthHandler.Root)
	app.Get("/health", healthHandler.Health)
	app.Get("/live", healthHandler.Live)
	app.Get("/ready", healthHandler.Ready)

	// Auth routes.
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Get("/me", middleware.AuthMiddleware(tokenSvc), authHandler.Me)

	// User CRUD routes (all require auth).
	users := app.Group("/api/v1/users", middleware.AuthMiddleware(tokenSvc))
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
	users.Post("/", userHandler.Create)
	users.Patch("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}
