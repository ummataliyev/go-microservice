package dto

import "time"

type CreateItemRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255" example:"Sourdough loaf"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000" example:"24-hour fermented; baked daily."`
	Quantity    int     `json:"quantity" validate:"gte=0" example:"12"`
}

type UpdateItemRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255" example:"Sourdough loaf v2"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000" example:"Updated recipe."`
	Quantity    *int    `json:"quantity,omitempty" validate:"omitempty,gte=0" example:"20"`
}

type ItemResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
