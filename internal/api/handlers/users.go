package handlers

import (
	"context"
	"strconv"

	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"

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

func (h *UserHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.svc.List(c.Context(), page, perPage)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

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

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user_id")
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

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid user_id")
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
