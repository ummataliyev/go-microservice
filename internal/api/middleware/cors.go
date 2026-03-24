package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thealish/go-microservice/internal/config"
)

// CORS returns a configured CORS middleware based on the application config.
func CORS(cfg config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization,X-Request-ID",
		AllowCredentials: true,
	})
}
