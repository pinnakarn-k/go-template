package httpresponse

import "transaction-api/internal/pagination"

type Response[T any] struct {
	Data T                `json:"data"`
	Meta *pagination.Meta `json:"meta,omitempty"`
}
