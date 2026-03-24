package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	domainerrors "go-microservice/internal/domain/errors"
	"go-microservice/internal/security"
)

// AuthMiddleware returns a Fiber handler that validates Bearer tokens using the
// provided TokenService. On success it stores the parsed claims in c.Locals("claims").
func AuthMiddleware(tokenSvc security.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			apiErr := domainerrors.NewUnauthorized("missing or invalid authorization header")
			return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenSvc.ValidateToken(tokenStr)
		if err != nil {
			apiErr := domainerrors.NewUnauthorized("invalid or expired token")
			return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
		}

		if claims.TokenType != "access" {
			apiErr := domainerrors.NewUnauthorized("invalid token type")
			return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
