package api

import (
	"strings"

	"go-microservice/internal/api/handlers"
	"go-microservice/internal/api/middleware"
	"go-microservice/internal/config"
	"go-microservice/internal/security"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

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
	app.Use(middleware.RequestID())
	app.Use(middleware.Timing(int64(cfg.Logging.SlowRequestThresholdMS)))
	app.Use(middleware.SecurityHeaders())

	var trustedHosts []string
	if cfg.Server.TrustedHosts != "" && cfg.Server.TrustedHosts != "*" {
		trustedHosts = strings.Split(cfg.Server.TrustedHosts, ",")
	}
	app.Use(middleware.TrustedHost(trustedHosts))

	app.Use(middleware.CORS(cfg.CORS))

	app.Get("/", healthHandler.Root)
	app.Get("/health", healthHandler.Health)
	app.Get("/live", healthHandler.Live)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/swagger/*", swagger.HandlerDefault)

	rl := middleware.NewRateLimiter(redisClient, cfg.RateLimit)
	app.Use(rl.Middleware())

	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Get("/me", middleware.AuthMiddleware(tokenSvc), authHandler.Me)

	users := app.Group("/api/v1/users", middleware.AuthMiddleware(tokenSvc))
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
	users.Post("/", userHandler.Create)
	users.Patch("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
}
