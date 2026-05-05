package handlers

import (
	"context"
	"strconv"

	"go-microservice/internal/dto"
	"go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
)

type UserServicer interface {
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.UserResponse], error)
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

type UserHandler struct {
	svc UserServicer
}

func NewUsers(svc UserServicer) *UserHandler {
	return &UserHandler{svc: svc}
}

// List godoc
// @Summary      List users (paginated)
// @Description  Returns a paginated list of active users.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        page     query int false "Page number (default 1)" default(1)
// @Param        per_page query int false "Items per page (default 20)" default(20)
// @Success      200 {object} dto.PaginatedResponse[dto.UserResponse] "Paginated user list"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Router       /api/v1/users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.svc.List(c.Context(), page, perPage)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Get godoc
// @Summary      Get a user by id
// @Description  Returns a single user. 404 if missing.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} dto.UserResponse "User record"
// @Failure      400 {object} errors.ErrorResponse "Invalid user_id"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Failure      404 {object} errors.ErrorResponse "User not found"
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid user_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Create godoc
// @Summary      Create a user
// @Description  Admin-style user creation. Hashes the password before storage.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateUserRequest true "User payload"
// @Success      201 {object} dto.UserResponse "Newly created user"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Failure      409 {object} errors.ErrorResponse "Email already exists"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Router       /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update godoc
// @Summary      Update a user (partial)
// @Description  Updates only the fields you send.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                    true  "User ID"
// @Param        request body dto.UpdateUserRequest  true  "Fields to update"
// @Success      200 {object} dto.UserResponse "Updated user"
// @Failure      400 {object} errors.ErrorResponse "Invalid user_id"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Failure      404 {object} errors.ErrorResponse "User not found"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Router       /api/v1/users/{id} [patch]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid user_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	var req dto.UpdateUserRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Delete godoc
// @Summary      Delete a user
// @Description  Soft-deletes the user record.
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Success      200 {object} dto.DeleteResponse "Delete confirmation"
// @Failure      400 {object} errors.ErrorResponse "Invalid user_id"
// @Failure      401 {object} errors.ErrorResponse "Missing or invalid token"
// @Failure      404 {object} errors.ErrorResponse "User not found"
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid user_id")
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

func parseIDParam(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
