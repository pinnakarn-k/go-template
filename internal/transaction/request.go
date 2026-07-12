package transaction

type SearchRequest struct {
	AccountID         int64  `json:"accountId" validate:"required"`
	TransactedAtStart string `json:"transactedAtStart" validate:"required"`
	TransactedAtEnd   string `json:"transactedAtEnd" validate:"required"`

	Page     int `json:"page" validate:"required,min=1"`
	PageSize int `json:"pageSize" validate:"required,min=1,max=100"`

	SortBy    string `json:"sortBy" validate:"omitempty,oneof=transacted_at amount fee status"`
	SortOrder string `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
}
