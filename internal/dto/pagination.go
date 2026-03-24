package dto

// PaginationRequest holds pagination query parameters.
type PaginationRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// NewPaginationRequest returns a PaginationRequest with sensible defaults.
func NewPaginationRequest(page, perPage int) PaginationRequest {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	return PaginationRequest{Page: page, PerPage: perPage}
}

// Offset returns the SQL-style offset for the current page.
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// PaginationMeta contains metadata about the paginated result set.
type PaginationMeta struct {
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	TotalItems  int  `json:"total_items"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// PaginatedResponse is a generic paginated response envelope.
type PaginatedResponse[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}
