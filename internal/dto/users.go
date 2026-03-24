package dto

import "time"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email" example:"new@example.com"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6" example:"newsecret123"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteResponse struct {
	Status string `json:"status"`
	ID     uint   `json:"id"`
}

type UserListResponse struct {
	Items []UserResponse `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
