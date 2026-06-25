package subscriptions

import "github.com/google/uuid"

type Subscription struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
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
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	TotalPages    int            `json:"total_pages"`
	HasPrev       bool           `json:"has_prev"`
	HasNext       bool           `json:"has_next"`
	PrevPage      int            `json:"prev_page"`
	NextPage      int            `json:"next_page"`
}
