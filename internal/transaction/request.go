package transaction

type SearchRequest struct {
	AccountID         int64  `json:"account_id" validate:"required"`
	TransactedAtStart string `json:"transacted_at_start" validate:"required"`
	TransactedAtEnd   string `json:"transacted_at_end" validate:"required"`

	Page     int `json:"page" validate:"required,min=1"`
	PageSize int `json:"page_size" validate:"required,min=1,max=100"`

	SortBy    string `json:"sort_by" validate:"omitempty,oneof=transacted_at amount fee status"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}
