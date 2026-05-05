package handlers

import (
	"context"

	"go-microservice/internal/dto"
	svcerrors "go-microservice/internal/errors"

	"github.com/gofiber/fiber/v2"
)

type ItemServicer interface {
	GetByID(ctx context.Context, id uint) (*dto.ItemResponse, error)
	List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.ItemResponse], error)
	Create(ctx context.Context, req dto.CreateItemRequest) (*dto.ItemResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateItemRequest) (*dto.ItemResponse, error)
	Delete(ctx context.Context, id uint) error
}

type ItemHandler struct {
	svc ItemServicer
}

func NewItems(svc ItemServicer) *ItemHandler {
	return &ItemHandler{svc: svc}
}

func (h *ItemHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.svc.List(c.Context(), page, perPage)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid item_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateItemRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid item_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	var req dto.UpdateItemRequest
	if err := validateBody(c, &req); err != nil {
		return err
	}

	result, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := svcerrors.NewBadRequest("invalid item_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
