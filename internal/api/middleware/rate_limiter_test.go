package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-microservice/internal/config"
	domainerrors "go-microservice/internal/errors"
)

func TestRateLimiter_SkipsWhenDisabled(t *testing.T) {
	cfg := config.RateLimitConfig{Enabled: false}
	rl := NewRateLimiter(nil, cfg)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRateLimiter_InMemoryFallback(t *testing.T) {
	cfg := config.RateLimitConfig{
		Enabled:  true,
		LimitGet: 3,
		TimeGet:  time.Minute,
		LimitPPD: 2,
		TimePPD:  time.Minute,
	}
	rl := NewRateLimiter(nil, cfg)

	app := fiber.New()
	app.Use(rl.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Requests 1-3 should pass (limit is 3).
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "request %d should pass", i+1)
	}

	// Request 4 should be rate limited.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Retry-After"))

	// Verify the error response body.
	var errResp domainerrors.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "TOO_MANY_REQUESTS", errResp.Error.Type)
}
