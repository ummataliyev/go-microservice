package middleware

import (
	"net/http"
	"strings"

	domainerrors "go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
)

func TrustedHost(trustedHosts []string) fiber.Handler {
	allowed := make(map[string]struct{}, len(trustedHosts))
	for _, h := range trustedHosts {
		allowed[strings.ToLower(strings.TrimSpace(h))] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		if len(allowed) == 0 {
			return c.Next()
		}

		host := c.Hostname()
		if idx := strings.LastIndex(host, ":"); idx != -1 {
			host = host[:idx]
		}
		host = strings.ToLower(host)

		if _, ok := allowed[host]; !ok {
			apiErr := &domainerrors.APIError{
				Type:       "MISDIRECTED_REQUEST",
				Message:    "untrusted host",
				StatusCode: http.StatusMisdirectedRequest,
			}
			return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
		}

		return c.Next()
	}
}
