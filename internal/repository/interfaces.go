package repository

import (
	"context"

	"go-microservice/internal/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}

type ItemRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Item, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Item, error)
	Create(ctx context.Context, item *models.Item) error
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
}
