package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/thealish/go-microservice/internal/config"
	domainerrors "github.com/thealish/go-microservice/internal/domain/errors"
)

// counter is used for the in-memory rate-limiting fallback.
type counter struct {
	count    int
	resetAt  time.Time
}

// RateLimiter implements a per-IP rate limiter backed by Redis (preferred) or an
// in-memory sync.Map fallback when Redis is unavailable.
type RateLimiter struct {
	redis    *redis.Client
	cfg      config.RateLimitConfig
	counters sync.Map
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(redisClient *redis.Client, cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
		cfg:   cfg,
	}
}

// Middleware returns a Fiber handler that enforces rate limits.
func (rl *RateLimiter) Middleware() fiber.Handler {
	if !rl.cfg.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		method := c.Method()

		var limit int
		var window time.Duration

		switch method {
		case fiber.MethodGet:
			limit = rl.cfg.LimitGet
			window = rl.cfg.TimeGet
		default: // POST, PUT, PATCH, DELETE, etc.
			limit = rl.cfg.LimitPPD
			window = rl.cfg.TimePPD
		}

		category := "read"
		if method != fiber.MethodGet {
			category = "write"
		}
		key := fmt.Sprintf("ratelimit:%s:%s", ip, category)

		var count int
		var err error

		if rl.redis != nil {
			count, err = rl.incrRedis(c.Context(), key, window)
		} else {
			count, err = rl.incrMemory(key, window)
		}
		if err != nil {
			// On error, allow the request through rather than blocking.
			return c.Next()
		}

		if count > limit {
			retryAfter := fmt.Sprintf("%d", int(window.Seconds()))
			c.Set("Retry-After", retryAfter)
			apiErr := domainerrors.NewTooManyRequests(retryAfter + "s")
			return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
		}

		return c.Next()
	}
}

func (rl *RateLimiter) incrRedis(ctx context.Context, key string, window time.Duration) (int, error) {
	pipe := rl.redis.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return int(incrCmd.Val()), nil
}

func (rl *RateLimiter) incrMemory(key string, window time.Duration) (int, error) {
	now := time.Now()
	val, loaded := rl.counters.Load(key)
	if loaded {
		c := val.(*counter)
		if now.After(c.resetAt) {
			c.count = 1
			c.resetAt = now.Add(window)
		} else {
			c.count++
		}
		return c.count, nil
	}

	c := &counter{
		count:   1,
		resetAt: now.Add(window),
	}
	actual, loaded := rl.counters.LoadOrStore(key, c)
	if loaded {
		existing := actual.(*counter)
		if now.After(existing.resetAt) {
			existing.count = 1
			existing.resetAt = now.Add(window)
		} else {
			existing.count++
		}
		return existing.count, nil
	}
	return c.count, nil
}
