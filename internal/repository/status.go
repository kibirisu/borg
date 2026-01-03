package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	Add(context.Context, db.AddStatusParams) error
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// Create implements StatusRepository.
func (r *statusRepository) Create(
	ctx context.Context,
	status db.CreateStatusParams,
) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}

// Add implements StatusRepository.
func (r *statusRepository) Add(ctx context.Context, status db.AddStatusParams) error {
	return r.q.AddStatus(ctx, status)
}
