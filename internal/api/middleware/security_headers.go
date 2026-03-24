package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeaders is a Fiber middleware that applies common security headers
// to every response and removes server-identifying headers.
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Remove server-identifying headers.
		c.Response().Header.Del("Server")
		c.Response().Header.Del("X-Powered-By")

		return err
	}
}
