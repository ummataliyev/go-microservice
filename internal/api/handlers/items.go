package handlers

import (
	"context"

	"go-microservice/internal/dto"
	"go-microservice/internal/errors"

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

// List godoc
// @Summary      List items (paginated)
// @Description  Returns a paginated list of items.
// @Tags         Items
// @Produce      json
// @Param        page     query int false "Page number (default 1)" default(1)
// @Param        per_page query int false "Items per page (default 20)" default(20)
// @Success      200 {object} dto.PaginatedResponse[dto.ItemResponse] "Paginated item list"
// @Router       /api/v1/items [get]
func (h *ItemHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)

	result, err := h.svc.List(c.Context(), page, perPage)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Get godoc
// @Summary      Get an item by id
// @Description  Returns a single item. 404 if missing.
// @Tags         Items
// @Produce      json
// @Param        id path int true "Item ID"
// @Success      200 {object} dto.ItemResponse "Item record"
// @Failure      400 {object} errors.ErrorResponse "Invalid item_id"
// @Failure      404 {object} errors.ErrorResponse "Item not found"
// @Router       /api/v1/items/{id} [get]
func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid item_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	result, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// Create godoc
// @Summary      Create an item
// @Description  Creates a new item.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateItemRequest true "Item payload"
// @Success      201 {object} dto.ItemResponse "Newly created item"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Router       /api/v1/items [post]
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

// Update godoc
// @Summary      Update an item (partial)
// @Description  Updates only the fields you send.
// @Tags         Items
// @Accept       json
// @Produce      json
// @Param        id      path int                    true  "Item ID"
// @Param        request body dto.UpdateItemRequest  true  "Fields to update"
// @Success      200 {object} dto.ItemResponse "Updated item"
// @Failure      400 {object} errors.ErrorResponse "Invalid item_id"
// @Failure      404 {object} errors.ErrorResponse "Item not found"
// @Failure      422 {object} errors.ErrorResponse "Validation error"
// @Router       /api/v1/items/{id} [patch]
func (h *ItemHandler) Update(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid item_id")
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

// Delete godoc
// @Summary      Delete an item
// @Description  Hard-deletes the item record.
// @Tags         Items
// @Produce      json
// @Param        id path int true "Item ID"
// @Success      204 "No content"
// @Failure      400 {object} errors.ErrorResponse "Invalid item_id"
// @Failure      404 {object} errors.ErrorResponse "Item not found"
// @Router       /api/v1/items/{id} [delete]
func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		apiErr := errors.NewBadRequest("invalid item_id")
		return c.Status(apiErr.StatusCode).JSON(apiErr.ToResponse())
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
