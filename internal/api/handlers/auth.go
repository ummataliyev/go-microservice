package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"
	"go-microservice/internal/security"
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
		return handleServiceError(c, err)
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
		return handleServiceError(c, err)
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
		return handleServiceError(c, err)
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
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
