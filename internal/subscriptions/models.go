package subscriptions

import "github.com/google/uuid"

type Subscription struct {
	ID          uuid.UUID `json:"id" swaggertype:"string" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"999"`
	UserID      uuid.UUID `json:"user_id" swaggertype:"string" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440001"`
	StartDate   string    `json:"start_date" example:"01-2024"`
	EndDate     *string   `json:"end_date,omitempty" example:"12-2024"`
}

type ListParams struct {
	ServiceName string
	UserID      uuid.UUID
	Sort        string
	Order       string
	Page        int
	PageSize    int
}

type ListResult struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int            `json:"total" example:"25"`
	Page          int            `json:"page" example:"1"`
	PageSize      int            `json:"page_size" example:"10"`
	TotalPages    int            `json:"total_pages" example:"3"`
	HasPrev       bool           `json:"has_prev" example:"false"`
	HasNext       bool           `json:"has_next" example:"true"`
	PrevPage      int            `json:"prev_page" example:"0"`
	NextPage      int            `json:"next_page" example:"2"`
}

type SumParams struct {
	StartDate   string
	EndDate     string
	UserID      uuid.UUID
	ServiceName string
}

type SumResult struct {
	Total int `json:"total" example:"1998"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}
