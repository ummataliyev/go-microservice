package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-microservice/internal/api/handlers"
	"go-microservice/internal/api/middleware"
	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserService) List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.UserResponse], error) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaginatedResponse[dto.UserResponse]), args.Error(1)
}

func (m *mockUserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserService) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestListUsers_Paginated(t *testing.T) {
	mockSvc := new(mockUserService)
	h := handlers.NewUsers(mockSvc)

	now := time.Now()
	paginatedResp := &dto.PaginatedResponse[dto.UserResponse]{
		Items: []dto.UserResponse{
			{ID: 1, Email: "user1@test.com", CreatedAt: now, UpdatedAt: now},
			{ID: 2, Email: "user2@test.com", CreatedAt: now, UpdatedAt: now},
		},
		Meta: dto.PaginationMeta{
			TotalPages:  1,
			CurrentPage: 1,
			TotalItems:  2,
			HasNext:     false,
			HasPrevious: false,
		},
	}
	mockSvc.On("List", mock.Anything, 1, 10).Return(paginatedResp, nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Get("/api/v1/users", h.List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=1&per_page=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.PaginatedResponse[dto.UserResponse]
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 1, result.Meta.CurrentPage)

	mockSvc.AssertExpectations(t)
}

func TestGetUser_Found(t *testing.T) {
	mockSvc := new(mockUserService)
	h := handlers.NewUsers(mockSvc)

	now := time.Now()
	user := &dto.UserResponse{ID: 1, Email: "user@test.com", CreatedAt: now, UpdatedAt: now}
	mockSvc.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Get("/api/v1/users/:id", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.UserResponse
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "user@test.com", result.Email)

	mockSvc.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockSvc := new(mockUserService)
	h := handlers.NewUsers(mockSvc)

	mockSvc.On("GetByID", mock.Anything, uint(999)).Return(nil, svcerrors.ErrUserNotFound)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Get("/api/v1/users/:id", h.Get)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockSvc.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h := handlers.NewUsers(mockSvc)

	now := time.Now()
	user := &dto.UserResponse{ID: 1, Email: "new@test.com", CreatedAt: now, UpdatedAt: now}
	mockSvc.On("Create", mock.Anything, dto.CreateUserRequest{
		Email:    "new@test.com",
		Password: "pass123",
	}).Return(user, nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Post("/api/v1/users", h.Create)

	body, _ := json.Marshal(map[string]string{
		"email":    "new@test.com",
		"password": "pass123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.UserResponse
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "new@test.com", result.Email)

	mockSvc.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	h := handlers.NewUsers(mockSvc)

	mockSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Delete("/api/v1/users/:id", h.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result dto.DeleteResponse
	require.NoError(t, json.Unmarshal(respBody, &result))
	assert.Equal(t, "success", result.Status)
	assert.Equal(t, uint(1), result.ID)

	mockSvc.AssertExpectations(t)
}
