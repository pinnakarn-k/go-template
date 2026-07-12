package transaction

import (
	"context"
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
