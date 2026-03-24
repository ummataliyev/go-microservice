package handlers

import (
	"context"

	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"
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
