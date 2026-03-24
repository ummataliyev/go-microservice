package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	repoerrors "go-microservice/internal/errors"
	"go-microservice/internal/models"
	"gorm.io/gorm"
)

// GORMUserRepository implements UserRepository using GORM.
// Works with any GORM-supported SQL database (currently PostgreSQL).
type GORMUserRepository struct {
	db *gorm.DB
}

// NewGORMUser returns a new GORMUserRepository.
func NewGORMUser(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{db: db}
}

func (r *GORMUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by id: %w", repoerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *GORMUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by email: %w", repoerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *GORMUserRepository) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	return users, nil
}

func (r *GORMUserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("create user: %w", repoerrors.ErrCannotCreate)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *GORMUserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", repoerrors.ErrCannotUpdate)
	}
	return nil
}

func (r *GORMUserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("delete user: %w", repoerrors.ErrCannotDelete)
	}
	return nil
}

func (r *GORMUserRepository) Restore(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("restore user: %w", err)
	}
	return nil
}

func (r *GORMUserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return count, nil
}

// isDuplicateKeyError checks whether the error is a duplicate-key / unique-constraint violation.
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "E11000")
}
