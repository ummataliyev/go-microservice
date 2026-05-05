package handlers

import (
	"context"

	"go-microservice/internal/dto"
	"go-microservice/internal/errors"
	"go-microservice/internal/security"

	"github.com/gofiber/fiber/v2"
)

type AuthServicer interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.TokenResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, clientIP string) (*dto.TokenResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.TokenResponse, error)
	GetCurrentUser(ctx context.Context, userID uint) (*dto.MeResponse, error)
}

type AuthHandler struct {
	svc AuthServicer
}

func NewAuth(svc AuthServicer) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns a JWT token pair.
// @Description  Email is normalised to lowercase before storage.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration payload"
// @Success      201 {object} dto.TokenResponse "Token pair for the newly created user"
// @Failure      409 {object} errors.ErrorResponse "Email already exists"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Failure      429 {object} errors.ErrorResponse "Rate limit exceeded"
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Register(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Login godoc
// @Summary      Log in with email and password
// @Description  Verifies credentials and returns a fresh access + refresh token pair.
// @Description  Email matching is case-insensitive. Repeated failures from the same IP/email
// @Description  are throttled (see AUTH_MAX_ATTEMPTS / AUTH_LOCKOUT_SECONDS).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login payload"
// @Success      200 {object} dto.TokenResponse "Fresh token pair"
// @Failure      401 {object} errors.ErrorResponse "Invalid credentials"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Failure      429 {object} errors.ErrorResponse "Account locked due to too many failed attempts or rate limit exceeded"
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Login(c.Context(), req, c.IP())
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Refresh godoc
// @Summary      Exchange a refresh token for a new pair
// @Description  Validates the refresh JWT, revokes it, and issues a fresh access + refresh pair.
// @Description  Reuse of a revoked refresh token returns 401.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh payload"
// @Success      200 {object} dto.TokenResponse "Fresh token pair"
// @Failure      401 {object} errors.ErrorResponse "Invalid or revoked refresh token"
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Refresh(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Me godoc
// @Summary      Return the current user
// @Description  Returns the authenticated user record. Requires a valid `Authorization: Bearer <access_token>` header.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.MeResponse "Authenticated user"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*security.Claims)
	if !ok || claims == nil {
		apiErr := errors.NewUnauthorized("missing or invalid claims")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetCurrentUser(c.Context(), claims.UserID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
