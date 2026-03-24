package dto

type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func NewPaginationRequest(page, perPage int) PaginationRequest {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return PaginationRequest{Page: page, PerPage: perPage}
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PerPage
}

type PaginationMeta struct {
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	TotalItems  int  `json:"total_items"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

type PaginatedResponse[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
