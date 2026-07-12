package transaction

import (
	"context"
	"errors"
	"time"
	"transaction-api/internal/pagination"
)

type Repository interface {
	Search(ctx context.Context, filter SearchFilter) ([]SearchRecord, int64, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

type SearchItem struct {
	ID              int64   `json:"id"`
	ReferenceCode   *string `json:"reference_code"`
	AccountID       *int64  `json:"account_id"`
	TransactionType *string `json:"transaction_type"`
	Side            *string `json:"side"`
	SideName        *string `json:"side_name"`
	Amount          *string `json:"amount"`
	Fee             *string `json:"fee"`
	Currency        *string `json:"currency"`
	Status          *string `json:"status"`
	Description     *string `json:"description"`
	TransactedAt    *string `json:"transacted_at"`
}

type SearchOutput struct {
	Items []SearchItem
	Meta  pagination.Meta
}

var ErrInvalidTimeRange = errors.New("transacted_at_start must not be after transacted_at_end")
var ErrInvalidDateFormat = errors.New("transacted_at_start and transacted_at_end must use YYYY-MM-DD format")

func (s *Service) Search(ctx context.Context, req SearchRequest) (SearchOutput, error) {
	transactedAtStart, err := time.Parse(
		"2006-01-02",
		req.TransactedAtStart,
	)
	if err != nil {
		return SearchOutput{}, ErrInvalidDateFormat
	}

	transactedAtEnd, err := time.Parse(
		"2006-01-02",
		req.TransactedAtEnd,
	)
	if err != nil {
		return SearchOutput{}, ErrInvalidDateFormat
	}

	if transactedAtStart.After(transactedAtEnd) {
		return SearchOutput{}, ErrInvalidTimeRange
	}

	page := pagination.Request{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	filter := SearchFilter{
		AccountID:         req.AccountID,
		TransactedAtStart: transactedAtStart,
		TransactedAtEnd:   transactedAtEnd,
		Limit:             page.Limit(),
		Offset:            page.Offset(),
		SortBy:            req.SortBy,
		SortOrder:         req.SortOrder,
	}

	records, total, err := s.repository.Search(ctx, filter)
	if err != nil {
		return SearchOutput{}, err
	}

	items := make([]SearchItem, 0, len(records))
	for _, record := range records {
		items = append(items, toSearchItem(record))
	}

	return SearchOutput{
		Items: items,
		Meta: pagination.NewMeta(
			req.Page,
			req.PageSize,
			total,
		),
	}, nil
}
