package service

import (
	"context"
	"errors"

	"go-microservice/internal/domain/dto"
	repoerrors "go-microservice/internal/domain/errors"
	"go-microservice/internal/domain/models"
	"go-microservice/internal/repository"
	"go-microservice/internal/security"
)

// UserService handles business logic for user operations.
type UserService struct {
	repo   repository.UserRepository
	hasher security.Hasher
}

// NewUsers creates a new UserService.
func NewUsers(repo repository.UserRepository, hasher security.Hasher) *UserService {
	return &UserService{repo: repo, hasher: hasher}
}

// GetByID fetches a user by ID and returns a UserResponse DTO.
func (s *UserService) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	resp := userToResponse(user)
	return &resp, nil
}

// List returns a paginated list of users.
func (s *UserService) List(ctx context.Context, page, perPage int) (*dto.PaginatedResponse[dto.UserResponse], error) {
	pag := dto.NewPaginationRequest(page, perPage)

	users, err := s.repo.GetAll(ctx, pag.PerPage, pag.Offset())
	if err != nil {
		return nil, err
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserResponse, len(users))
	for i := range users {
		items[i] = userToResponse(&users[i])
	}

	totalItems := int(count)
	totalPages := totalItems / pag.PerPage
	if totalItems%pag.PerPage != 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.UserResponse]{
		Items: items,
		Meta: dto.PaginationMeta{
			TotalPages:  totalPages,
			CurrentPage: pag.Page,
			TotalItems:  totalItems,
			HasNext:     pag.Page < totalPages,
			HasPrevious: pag.Page > 1,
		},
	}, nil
}

// Create hashes the password, creates the user, and returns a UserResponse DTO.
func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashed, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:          req.Email,
		HashedPassword: hashed,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, repoerrors.ErrCannotCreate) {
			return nil, repoerrors.ErrUserAlreadyExists
		}
		return nil, err
	}

	resp := userToResponse(user)
	return &resp, nil
}

// Update updates an existing user's fields and returns the updated UserResponse DTO.
func (s *UserService) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashed, err := s.hasher.Hash(*req.Password)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = hashed
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := userToResponse(user)
	return &resp, nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// userToResponse converts a User model to a UserResponse DTO.
func userToResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
