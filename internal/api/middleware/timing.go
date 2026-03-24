package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Timing(slowThresholdMS int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		durationMS := float64(time.Since(start).Nanoseconds()) / 1e6
		c.Set("X-Process-Time", fmt.Sprintf("%.2f", durationMS))

		if int64(durationMS) > slowThresholdMS {
			log.Warn().
				Str("path", c.Path()).
				Str("method", c.Method()).
				Float64("duration_ms", durationMS).
				Msg("slow request detected")
		}

		return err
	}
}
