package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thealish/go-microservice/internal/api/handlers"
)

func TestHealthEndpoint(t *testing.T) {
	app := fiber.New()
	h := handlers.NewHealth("test-app", "1.0.0")
	app.Get("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "ok", result["status"])
}

func TestLiveEndpoint(t *testing.T) {
	app := fiber.New()
	h := handlers.NewHealth("test-app", "1.0.0")
	app.Get("/live", h.Live)

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "alive", result["status"])
}

func TestReadyEndpoint(t *testing.T) {
	app := fiber.New()
	h := handlers.NewHealth("test-app", "1.0.0")
	app.Get("/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "ready", result["status"])
}

func TestRootEndpoint(t *testing.T) {
	app := fiber.New()
	h := handlers.NewHealth("test-app", "1.0.0")
	app.Get("/", h.Root)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	assert.Equal(t, "test-app", result["app"])
	assert.Equal(t, "1.0.0", result["version"])
	assert.Equal(t, "running", result["status"])
}
