package repository

import (
	"context"

	"go-microservice/internal/domain/models"
)

// UserRepository defines the contract for user persistence operations.
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
