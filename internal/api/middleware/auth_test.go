package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thealish/go-microservice/internal/domain/errors"
	"github.com/thealish/go-microservice/internal/security"
)

// mockTokenService is a test double for security.TokenService.
type mockTokenService struct {
	claims *security.Claims
	err    error
}

func (m *mockTokenService) GenerateAccessToken(userID uint, email string) (string, error) {
	return "access-token", nil
}

func (m *mockTokenService) GenerateRefreshToken(userID uint, email string) (string, error) {
	return "refresh-token", nil
}

func (m *mockTokenService) ValidateToken(tokenString string) (*security.Claims, error) {
	return m.claims, m.err
}

func TestAuth_ValidToken(t *testing.T) {
	svc := &mockTokenService{
		claims: &security.Claims{
			UserID:    1,
			Email:     "user@example.com",
			TokenType: "access",
		},
	}

	app := fiber.New()
	app.Use(AuthMiddleware(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*security.Claims)
		return c.SendString(fmt.Sprintf("user:%d", claims.UserID))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_MissingHeader(t *testing.T) {
	svc := &mockTokenService{}

	app := fiber.New()
	app.Use(AuthMiddleware(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errResp errors.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", errResp.Error.Type)
}

func TestAuth_InvalidToken(t *testing.T) {
	svc := &mockTokenService{
		err: fmt.Errorf("token expired"),
	}

	app := fiber.New()
	app.Use(AuthMiddleware(svc))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errResp errors.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", errResp.Error.Type)
}
