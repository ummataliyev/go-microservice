package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-microservice/internal/config"
	domainerrors "go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type counter struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	redis    *redis.Client
	cfg      config.RateLimitConfig
	counters sync.Map
}

func NewRateLimiter(redisClient *redis.Client, cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
		cfg:   cfg,
	}
}

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
		default:
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
			count = rl.incrMemory(key, window)
		}
		if err != nil {
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

func (rl *RateLimiter) incrMemory(key string, window time.Duration) int {
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
		return c.count
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
		return existing.count
	}
	return c.count
}
