package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/the_elita/go-microservice/internal/domain/dto"
	svcerrors "github.com/the_elita/go-microservice/internal/domain/errors"
	"github.com/the_elita/go-microservice/internal/security"
)

// AuthServicer defines the interface for the auth service, consumed by the handler.
type AuthServicer interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, clientIP string) (*dto.TokenResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.TokenResponse, error)
	GetCurrentUser(ctx context.Context, userID uint) (*dto.MeResponse, error)
}

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	svc AuthServicer
}

// NewAuth creates a new AuthHandler.
func NewAuth(svc AuthServicer) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles user registration.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := svcerrors.NewBadRequest("invalid request body")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if req.Email == "" || req.Password == "" {
		apiErr := svcerrors.NewBadRequest("email and password are required")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Login handles user login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := svcerrors.NewBadRequest("invalid request body")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if req.Email == "" || req.Password == "" {
		apiErr := svcerrors.NewBadRequest("email and password are required")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.Login(c.Context(), req, c.IP())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Refresh handles token refresh.
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := svcerrors.NewBadRequest("invalid request body")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if req.RefreshToken == "" {
		apiErr := svcerrors.NewBadRequest("refresh_token is required")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.Refresh(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Me returns the current authenticated user's profile.
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*security.Claims)
	if !ok || claims == nil {
		apiErr := svcerrors.NewUnauthorized("missing or invalid claims")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetCurrentUser(c.Context(), claims.UserID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// handleError maps service-layer errors to appropriate HTTP responses.
func (h *AuthHandler) handleError(c *fiber.Ctx, err error) error {
	return handleServiceError(c, err)
}

// handleServiceError is the shared error mapper used by all handlers.
func handleServiceError(c *fiber.Ctx, err error) error {
	var apiErr *svcerrors.APIError

	switch {
	case errors.Is(err, svcerrors.ErrUserAlreadyExists):
		apiErr = svcerrors.NewConflict(err.Error())
	case errors.Is(err, svcerrors.ErrInvalidCredentials):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrLoginLocked):
		apiErr = svcerrors.NewTooManyRequests("30m")
	case errors.Is(err, svcerrors.ErrInvalidToken):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrInvalidTokenType):
		apiErr = svcerrors.NewUnauthorized(err.Error())
	case errors.Is(err, svcerrors.ErrUserNotFound):
		apiErr = svcerrors.NewNotFound(err.Error())
	default:
		apiErr = svcerrors.NewInternal("internal server error")
	}

	return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
}
