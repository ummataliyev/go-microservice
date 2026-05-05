package repository

import (
	"context"
	"errors"
	"fmt"

	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"

	"gorm.io/gorm"
)

type GORMItemRepository struct {
	db *gorm.DB
}

func NewGORMItem(db *gorm.DB) *GORMItemRepository {
	return &GORMItemRepository{db: db}
}

func (r *GORMItemRepository) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get item by id: %w", repoerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get item by id: %w", err)
	}
	return &item, nil
}

func (r *GORMItemRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get all items: %w", err)
	}
	return items, nil
}

func (r *GORMItemRepository) Create(ctx context.Context, item *models.Item) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create item: %w", err)
	}
	return nil
}

func (r *GORMItemRepository) Update(ctx context.Context, item *models.Item) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		return fmt.Errorf("update item: %w", repoerrors.ErrCannotUpdate)
	}
	return nil
}

func (r *GORMItemRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Item{}, id).Error; err != nil {
		return fmt.Errorf("delete item: %w", repoerrors.ErrCannotDelete)
	}
	return nil
}

func (r *GORMItemRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count items: %w", err)
	}
	return count, nil
}
