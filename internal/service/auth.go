package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go-microservice/internal/config"
	"go-microservice/internal/dto"
	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"
	"go-microservice/internal/repository"
	"go-microservice/internal/security"
)

type AuthService struct {
	repo     repository.UserRepository
	tokenSvc security.TokenService
	hasher   security.Hasher
	redis    *redis.Client
	cfg      config.AuthConfig
}

func NewAuth(
	repo repository.UserRepository,
	tokenSvc security.TokenService,
	hasher security.Hasher,
	redisClient *redis.Client,
	cfg config.AuthConfig,
) *AuthService {
	return &AuthService{
		repo:     repo,
		tokenSvc: tokenSvc,
		hasher:   hasher,
		redis:    redisClient,
		cfg:      cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error) {
	// Check if email already exists.
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, repoerrors.ErrUserAlreadyExists
	}

	hashed, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          req.Email,
		HashedPassword: hashed,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, repoerrors.ErrCannotCreate) {
			return nil, repoerrors.ErrUserAlreadyExists
		}
		return nil, err
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	resp := dto.NewTokenResponse(accessToken, refreshToken)
	return &resp, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, clientIP string) (*dto.TokenResponse, error) {
	if err := s.checkLoginLock(ctx, req.Email, clientIP); err != nil {
		return nil, err
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.hasher.Verify(req.Password, user.HashedPassword); err != nil {
		s.trackFailedAttempt(ctx, req.Email, clientIP)
		return nil, repoerrors.ErrInvalidCredentials
	}

	s.clearFailedAttempts(ctx, req.Email, clientIP)

	accessToken, err := s.tokenSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenSvc.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	resp := dto.NewTokenResponse(accessToken, refreshToken)
	return &resp, nil
}

func (s *AuthService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.TokenResponse, error) {
	claims, err := s.tokenSvc.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, repoerrors.ErrInvalidToken
	}

	if claims.TokenType != "refresh" {
		return nil, repoerrors.ErrInvalidTokenType
	}

	accessToken, err := s.tokenSvc.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	resp := dto.NewTokenResponse(accessToken, req.RefreshToken)
	return &resp, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*dto.MeResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	return &dto.MeResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) checkLoginLock(ctx context.Context, email, ip string) error {
	if s.redis == nil {
		return nil
	}

	key := fmt.Sprintf("login_attempts:%s:%s", email, ip)
	val, err := s.redis.Get(ctx, key).Int()
	if err != nil {
		return nil //nolint:nilerr // redis key not found means no lock
	}

	if val >= s.cfg.MaxAttempts {
		return repoerrors.ErrLoginLocked
	}

	return nil
}

func (s *AuthService) trackFailedAttempt(ctx context.Context, email, ip string) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("login_attempts:%s:%s", email, ip)
	s.redis.Incr(ctx, key)
	s.redis.Expire(ctx, key, time.Duration(s.cfg.WindowSeconds)*time.Second)
}

func (s *AuthService) clearFailedAttempts(ctx context.Context, email, ip string) {
	if s.redis == nil {
		return
	}

	key := fmt.Sprintf("login_attempts:%s:%s", email, ip)
	s.redis.Del(ctx, key)
}
