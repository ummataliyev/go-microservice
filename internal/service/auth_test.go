package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/thealish/go-microservice/internal/config"
	"github.com/thealish/go-microservice/internal/domain/dto"
	svcerrors "github.com/thealish/go-microservice/internal/domain/errors"
	"github.com/thealish/go-microservice/internal/security"
)

// --- Mock TokenService ---

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(tokenString string) (*security.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*security.Claims), args.Error(1)
}

// --- Helpers ---

func defaultAuthCfg() config.AuthConfig {
	return config.AuthConfig{
		MaxAttempts:    5,
		WindowSeconds:  900,
		LockoutSeconds: 1800,
	}
}

func newAuthService() (*AuthService, *MockUserRepository, *MockHasher, *MockTokenService) {
	repo := new(MockUserRepository)
	hasher := new(MockHasher)
	tokenSvc := new(MockTokenService)
	svc := NewAuth(repo, tokenSvc, hasher, nil, defaultAuthCfg())
	return svc, repo, hasher, tokenSvc
}

// --- Register Tests ---

func TestRegister_Success(t *testing.T) {
	svc, repo, hasher, tokenSvc := newAuthService()

	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, svcerrors.ErrNotFound)
	hasher.On("Hash", "password123").Return("hashed_pw", nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
	tokenSvc.On("GenerateAccessToken", mock.AnythingOfType("uint"), "test@example.com").Return("access_token", nil)
	tokenSvc.On("GenerateRefreshToken", mock.AnythingOfType("uint"), "test@example.com").Return("refresh_token", nil)

	req := dto.RegisterRequest{Email: "test@example.com", Password: "password123"}
	resp, err := svc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
	assert.Equal(t, "bearer", resp.TokenType)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc, repo, _, _ := newAuthService()

	user := sampleUser()
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	req := dto.RegisterRequest{Email: "test@example.com", Password: "password123"}
	resp, err := svc.Register(context.Background(), req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, svcerrors.ErrUserAlreadyExists)
	repo.AssertExpectations(t)
}

// --- Login Tests ---

func TestLogin_Success(t *testing.T) {
	svc, repo, hasher, tokenSvc := newAuthService()

	user := sampleUser()
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	hasher.On("Verify", "password123", "hashed_pw").Return(nil)
	tokenSvc.On("GenerateAccessToken", user.ID, user.Email).Return("access_token", nil)
	tokenSvc.On("GenerateRefreshToken", user.ID, user.Email).Return("refresh_token", nil)

	req := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
	resp, err := svc.Login(context.Background(), req, "127.0.0.1")

	require.NoError(t, err)
	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
	repo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, repo, hasher, _ := newAuthService()

	user := sampleUser()
	repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	hasher.On("Verify", "wrong_password", "hashed_pw").Return(errors.New("mismatch"))

	req := dto.LoginRequest{Email: "test@example.com", Password: "wrong_password"}
	resp, err := svc.Login(context.Background(), req, "127.0.0.1")

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, svcerrors.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, repo, _, _ := newAuthService()

	repo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, svcerrors.ErrNotFound)

	req := dto.LoginRequest{Email: "notfound@example.com", Password: "password123"}
	resp, err := svc.Login(context.Background(), req, "127.0.0.1")

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, svcerrors.ErrInvalidCredentials)
	repo.AssertExpectations(t)
}

// --- Refresh Tests ---

func TestRefresh_Success(t *testing.T) {
	svc, _, _, tokenSvc := newAuthService()

	claims := &security.Claims{
		UserID:    1,
		Email:     "test@example.com",
		TokenType: "refresh",
	}
	tokenSvc.On("ValidateToken", "valid_refresh_token").Return(claims, nil)
	tokenSvc.On("GenerateAccessToken", uint(1), "test@example.com").Return("new_access_token", nil)

	req := dto.RefreshRequest{RefreshToken: "valid_refresh_token"}
	resp, err := svc.Refresh(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "new_access_token", resp.AccessToken)
	assert.Equal(t, "valid_refresh_token", resp.RefreshToken)
	tokenSvc.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	svc, _, _, tokenSvc := newAuthService()

	tokenSvc.On("ValidateToken", "bad_token").Return(nil, errors.New("invalid"))

	req := dto.RefreshRequest{RefreshToken: "bad_token"}
	resp, err := svc.Refresh(context.Background(), req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, svcerrors.ErrInvalidToken)
	tokenSvc.AssertExpectations(t)
}

func TestRefresh_WrongTokenType(t *testing.T) {
	svc, _, _, tokenSvc := newAuthService()

	claims := &security.Claims{
		UserID:    1,
		Email:     "test@example.com",
		TokenType: "access",
	}
	tokenSvc.On("ValidateToken", "access_token_used_as_refresh").Return(claims, nil)

	req := dto.RefreshRequest{RefreshToken: "access_token_used_as_refresh"}
	resp, err := svc.Refresh(context.Background(), req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, svcerrors.ErrInvalidTokenType)
	tokenSvc.AssertExpectations(t)
}

// --- GetCurrentUser Tests ---

func TestGetCurrentUser(t *testing.T) {
	svc, repo, _, _ := newAuthService()

	user := sampleUser()
	repo.On("GetByID", mock.Anything, uint(1)).Return(user, nil)

	resp, err := svc.GetCurrentUser(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.CreatedAt, resp.CreatedAt)
	repo.AssertExpectations(t)
}
