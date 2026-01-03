package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	GetById(context.Context, int) (db.Status, error)
	GetByIdWithMetadata(context.Context, int) (db.GetStatusByIdWithMetadataRow, error)
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// Create implements StatusRepository.
func (r *statusRepository) Create(ctx context.Context, status db.CreateStatusParams) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}
// GetById implements StatusRepository.
func (r *statusRepository) GetById(ctx context.Context, id int) (db.Status, error) {
	return r.q.GetStatusById(ctx, int32(id))
}
// GetById implements StatusRepository.
func (r *statusRepository) GetByIdWithMetadata(ctx context.Context, id int) (db.GetStatusByIdWithMetadataRow, error) {
	return r.q.GetStatusByIdWithMetadata(ctx, int32(id))
}
