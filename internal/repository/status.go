package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) error
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// Create implements StatusRepository.
func (r *statusRepository) Create(ctx context.Context, status db.CreateStatusParams) error {
	return r.q.CreateStatus(ctx, status)
}
