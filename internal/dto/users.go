package dto

import "time"

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest is the payload for updating a user.
// Pointer fields are optional — nil means "do not change".
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

// UserResponse is the standard user representation returned by user endpoints.
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteResponse confirms a successful deletion.
type DeleteResponse struct {
	Status string `json:"status"`
	ID     uint   `json:"id"`
}
