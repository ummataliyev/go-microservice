package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	domainerrors "go-microservice/internal/errors"
)

func TestTrustedHost_AllowsValidHost(t *testing.T) {
	app := fiber.New()
	app.Use(TrustedHost([]string{"example.com", "api.example.com"}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTrustedHost_RejectsInvalidHost(t *testing.T) {
	app := fiber.New()
	app.Use(TrustedHost([]string{"example.com"}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "http://evil.com/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusMisdirectedRequest, resp.StatusCode)

	var errResp domainerrors.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "MISDIRECTED_REQUEST", errResp.Error.Type)
}

func TestTrustedHost_AllowsAllWhenEmpty(t *testing.T) {
	app := fiber.New()
	app.Use(TrustedHost([]string{}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "http://anything.com/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
