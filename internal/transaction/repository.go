package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type SearchFilter struct {
	AccountID         int64
	TransactedAtStart time.Time
	TransactedAtEnd   time.Time
	Limit             int
	Offset            int
}

type SearchRecord struct {
	ID              int64
	ReferenceCode   *string
	AccountID       *int64
	TransactionType *string
	Side            *string
	Amount          *decimal.Decimal
	Fee             *decimal.Decimal
	Currency        *string
	Status          *string
	Description     *string
	TransactedAt    *time.Time
}

const searchBase = `
FROM transactions
WHERE account_id = $1
  AND transacted_at >= $2
  AND transacted_at <= $3
`

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Search(
	ctx context.Context,
	filter SearchFilter,
) ([]SearchRecord, int64, error) {
	total, err := r.count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	records, err := r.search(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *PostgresRepository) count(
	ctx context.Context,
	filter SearchFilter,
) (int64, error) {
	const query = `
		SELECT COUNT(*)
	` + searchBase

	var total int64

	err := r.db.QueryRow(
		ctx,
		query,
		filter.AccountID,
		filter.TransactedAtStart,
		filter.TransactedAtEnd,
	).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count transactions: %w", err)
	}

	return total, nil
}

func (r *PostgresRepository) search(
	ctx context.Context,
	filter SearchFilter,
) ([]SearchRecord, error) {
	const query = `
		SELECT
			id,
			reference_code,
			account_id,
			transaction_type,
			side,
			amount,
			fee,
			currency,
			status,
			description,
			transacted_at
	` + searchBase + `
		ORDER BY transacted_at DESC
		LIMIT $4
		OFFSET $5
	`

	rows, err := r.db.Query(
		ctx,
		query,
		filter.AccountID,
		filter.TransactedAtStart,
		filter.TransactedAtEnd,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("search transactions: %w", err)
	}
	defer rows.Close()

	records := make([]SearchRecord, 0)
	for rows.Next() {
		var record SearchRecord

		if err := rows.Scan(
			&record.ID,
			&record.ReferenceCode,
			&record.AccountID,
			&record.TransactionType,
			&record.Side,
			&record.Amount,
			&record.Fee,
			&record.Currency,
			&record.Status,
			&record.Description,
			&record.TransactedAt,
		); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transactions: %w", err)
	}

	return records, nil
}
