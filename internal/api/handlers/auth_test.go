package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-microservice/internal/api/handlers"
	"go-microservice/internal/api/middleware"
	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, req dto.LoginRequest, clientIP string) (*dto.TokenResponse, error) {
	args := m.Called(ctx, req, clientIP)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *mockAuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.TokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *mockAuthService) GetCurrentUser(ctx context.Context, userID uint) (*dto.MeResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MeResponse), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockSvc := new(mockAuthService)
	h := handlers.NewAuth(mockSvc)

	tokenResp := &dto.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		TokenType:    "bearer",
	}
	mockSvc.On("Register", mock.Anything, dto.RegisterRequest{
		Email:    "test@test.com",
		Password: "pass123",
	}).Return(tokenResp, nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Post("/api/v1/auth/register", h.Register)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@test.com",
		"password": "pass123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.TokenResponse
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.Equal(t, "bearer", result.TokenType)

	mockSvc.AssertExpectations(t)
}

func TestRegister_InvalidBody(t *testing.T) {
	mockSvc := new(mockAuthService)
	h := handlers.NewAuth(mockSvc)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Post("/api/v1/auth/register", h.Register)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	mockSvc := new(mockAuthService)
	h := handlers.NewAuth(mockSvc)

	tokenResp := &dto.TokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		TokenType:    "bearer",
	}
	mockSvc.On("Login", mock.Anything, dto.LoginRequest{
		Email:    "test@test.com",
		Password: "pass123",
	}, "0.0.0.0").Return(tokenResp, nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Post("/api/v1/auth/login", h.Login)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@test.com",
		"password": "pass123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.TokenResponse
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)

	mockSvc.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockSvc := new(mockAuthService)
	h := handlers.NewAuth(mockSvc)

	mockSvc.On("Login", mock.Anything, dto.LoginRequest{
		Email:    "test@test.com",
		Password: "wrong",
	}, "0.0.0.0").Return(nil, svcerrors.ErrInvalidCredentials)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Post("/api/v1/auth/login", h.Login)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@test.com",
		"password": "wrong",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockSvc.AssertExpectations(t)
}
