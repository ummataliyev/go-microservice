package service

import (
	"context"
	"errors"

	"go-microservice/internal/dto"
	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"
	"go-microservice/internal/repository"
)

type ItemService struct {
	repo repository.ItemRepository
}

func NewItems(repo repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetByID(ctx context.Context, id uint) (*dto.ItemResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrItemNotFound
		}
		return nil, err
	}
	resp := itemToResponse(item)
	return &resp, nil
}

func (s *ItemService) List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.ItemResponse], error) {
	pag := dto.NewPaginationRequest(page, perPage)

	items, err := s.repo.GetAll(ctx, pag.PerPage, pag.Offset())
	if err != nil {
		return nil, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	rows := make([]dto.ItemResponse, len(items))
	for i := range items {
		rows[i] = itemToResponse(&items[i])
	}

	totalItems := int(count)
	totalPages := totalItems / pag.PerPage
	if totalItems%pag.PerPage != 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.ItemResponse]{
		Items: rows,
		Meta: dto.PaginationMeta{
			TotalPages:  totalPages,
			CurrentPage: pag.Page,
			TotalItems:  totalItems,
			HasNext:     pag.Page < totalPages,
			HasPrevious: pag.Page > 1,
		},
	}, nil
}

func (s *ItemService) Create(ctx context.Context, req dto.CreateItemRequest) (*dto.ItemResponse, error) {
	item := &models.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	resp := itemToResponse(item)
	return &resp, nil
}

func (s *ItemService) Update(ctx context.Context, id uint, req dto.UpdateItemRequest) (*dto.ItemResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrItemNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	resp := itemToResponse(item)
	return &resp, nil
}

func (s *ItemService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return repoerrors.ErrItemNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func itemToResponse(item *models.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Quantity:    item.Quantity,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
