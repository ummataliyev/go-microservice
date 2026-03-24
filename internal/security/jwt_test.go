package security

import (
	"testing"
	"time"

	"go-microservice/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestJWTService() *JWTService {
	cfg := config.JWTConfig{
		SecretKey:          "test-secret-key-for-unit-tests",
		Algorithm:          "HS256",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	return NewJWTService(cfg)
}

func TestGenerateAccessToken_ValidClaims(t *testing.T) {
	svc := newTestJWTService()
	userID := uint(42)
	email := "user@example.com"

	tokenStr, err := svc.GenerateAccessToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := svc.ValidateToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "access", claims.TokenType)
}

func TestGenerateRefreshToken_ValidClaims(t *testing.T) {
	svc := newTestJWTService()
	userID := uint(7)
	email := "refresh@example.com"

	tokenStr, err := svc.GenerateRefreshToken(userID, email)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := svc.ValidateToken(tokenStr)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "refresh", claims.TokenType)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	cfg := config.JWTConfig{
		SecretKey:          "test-secret-key-for-unit-tests",
		Algorithm:          "HS256",
		AccessTokenExpiry:  -1 * time.Second, // already expired
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	svc := NewJWTService(cfg)

	tokenStr, err := svc.GenerateAccessToken(1, "expired@example.com")
	require.NoError(t, err)

	_, err = svc.ValidateToken(tokenStr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestValidateToken_WrongTokenType(t *testing.T) {
	svc := newTestJWTService()

	tokenStr, err := svc.GenerateRefreshToken(1, "test@example.com")
	require.NoError(t, err)

	claims, err := svc.ValidateToken(tokenStr)
	require.NoError(t, err)
	assert.NotEqual(t, "access", claims.TokenType)
	assert.Equal(t, "refresh", claims.TokenType)
}
