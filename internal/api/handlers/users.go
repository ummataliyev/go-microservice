package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"
)

// UserServicer defines the interface for the user service, consumed by the handler.
type UserServicer interface {
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.UserResponse], error)
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

// UserHandler handles user CRUD HTTP requests.
type UserHandler struct {
	svc UserServicer
}

// NewUsers creates a new UserHandler.
func NewUsers(svc UserServicer) *UserHandler {
	return &UserHandler{svc: svc}
}

// List returns a paginated list of users.
func (h *UserHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.svc.List(c.Context(), page, perPage)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Get returns a single user by ID.
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Create creates a new user.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := svcerrors.NewBadRequest("invalid request body")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if req.Email == "" || req.Password == "" {
		apiErr := svcerrors.NewBadRequest("email and password are required")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update updates an existing user.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		apiErr := svcerrors.NewBadRequest("invalid request body")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Delete removes a user by ID.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.DeleteResponse{
		Status: "success",
		ID:     id,
	})
}

// parseIDParam extracts and validates the :id route parameter.
func parseIDParam(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
