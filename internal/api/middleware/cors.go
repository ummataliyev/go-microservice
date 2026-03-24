package middleware

import (
	"go-microservice/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS(cfg config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization,X-Request-ID",
		AllowCredentials: true,
	})
}
