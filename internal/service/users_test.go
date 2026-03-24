package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go-microservice/internal/dto"
	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"
)

// --- Mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Restore(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Verify(password, hash string) error {
	args := m.Called(password, hash)
	return args.Error(0)
}

// --- Tests ---

func sampleUser() *models.User {
	return &models.User{
		ID:             1,
		Email:          "test@example.com",
		HashedPassword: "hashed_pw",
		CreatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestGetUser_Found(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	user := sampleUser()
	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	resp, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.CreatedAt, resp.CreatedAt)
	repo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	repo.On("GetByID", mock.Anything, uint(99)).Return(nil, repoerrors.ErrNotFound)

	resp, err := svc.GetByID(context.Background(), 99)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, repoerrors.ErrUserNotFound)
	repo.AssertExpectations(t)
}

func TestListUsers_WithPagination(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	users := []models.User{*sampleUser()}
	// page=2, perPage=10 => offset=10, limit=10
	repo.On("GetAll", mock.Anything, 10, 10).Return(users, nil)
	repo.On("Count", mock.Anything).Return(int64(15), nil)

	resp, err := svc.List(context.Background(), 2, 10)

	require.NoError(t, err)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, 15, resp.Meta.TotalItems)
	assert.Equal(t, 2, resp.Meta.TotalPages)
	assert.Equal(t, 2, resp.Meta.CurrentPage)
	assert.False(t, resp.Meta.HasNext)
	assert.True(t, resp.Meta.HasPrevious)
	repo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	hasher.On("Hash", "password123").Return("hashed_pw", nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := svc.Create(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "test@example.com", resp.Email)
	hasher.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	hasher.On("Hash", "password123").Return("hashed_pw", nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(repoerrors.ErrCannotCreate)

	req := dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := svc.Create(context.Background(), req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, repoerrors.ErrUserAlreadyExists)
	repo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	svc := NewUsers(repo, hasher)

	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
