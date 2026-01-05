package repository

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	GetById(context.Context, int) (db.Status, error)
	GetShares(context.Context, int) ([]db.Status, error)
	GetLocalStatuses(context.Context) ([]db.GetLocalStatusesRow, error)
	GetByIdWithMetadata(context.Context, int) (db.GetStatusByIdWithMetadataRow, error)
	GetComments(context.Context, int) ([]db.Status, error)
	Update(context.Context, db.UpdateStatusParams) (db.Status, error)
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
// GetShares implements StatusRepository.
func (r *statusRepository) GetShares(ctx context.Context, id int) ([]db.Status, error) {
	return r.q.GetStatusShares(ctx, sql.NullInt32{Int32: int32(id), Valid: true})
}
// GetLocalStatuses implements StatusRepository.
func (r *statusRepository) GetLocalStatuses(ctx context.Context) ([]db.GetLocalStatusesRow , error) {
	return r.q.GetLocalStatuses(ctx)
}
// GetComments implements StatusRepository.
func (r *statusRepository) GetComments(ctx context.Context, id int) ([]db.Status, error) {
    return r.q.GetStatusComments(ctx, int32(id))
}
// Update implements StatusRepository.
func (r *statusRepository) Update(ctx context.Context, params db.UpdateStatusParams) (db.Status, error) {
    return r.q.UpdateStatus(ctx, params)
}