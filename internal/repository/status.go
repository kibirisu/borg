package repository

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	Add(context.Context, db.AddStatusParams) error
	GetByID(context.Context, int) (db.Status, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetShares(context.Context, int) ([]db.Status, error)
	GetLocalStatuses(context.Context) ([]db.GetLocalStatusesRow, error)
	GetByIDWithMetadata(context.Context, int) (db.GetStatusByIdWithMetadataRow, error)
	GetSharedPostsByAccountId(context.Context, int) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountId(context.Context, int) ([]db.GetTimelinePostsByAccountIdRow, error)
	DeleteByID(context.Context, int32) error
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

// GetById implements StatusRepository.
func (r *statusRepository) GetByID(ctx context.Context, id int) (db.Status, error) {
	return r.q.GetStatusById(ctx, int32(id))
}

// GetByURI implements StatusRepository.
func (r *statusRepository) GetByURI(ctx context.Context, uri string) (db.Status, error) {
	return r.q.GetStatusByURI(ctx, uri)
}

// GetById implements StatusRepository.
func (r *statusRepository) GetByIDWithMetadata(
	ctx context.Context,
	id int,
) (db.GetStatusByIdWithMetadataRow, error) {
	return r.q.GetStatusByIdWithMetadata(ctx, int32(id))
}

// GetShares implements StatusRepository.
func (r *statusRepository) GetShares(ctx context.Context, id int) ([]db.Status, error) {
	return r.q.GetStatusShares(ctx, sql.NullInt32{Int32: int32(id), Valid: true})
}

// GetLocalStatuses implements StatusRepository.
func (r *statusRepository) GetLocalStatuses(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return r.q.GetLocalStatuses(ctx)
}

// GetSharedPostsByAccountId implements StatusRepository.
func (r *statusRepository) GetSharedPostsByAccountId(
	ctx context.Context,
	accountID int,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	return r.q.GetSharedPostsByAccountId(ctx, int32(accountID))
}

// GetTimelinePostsByAccountId implements StatusRepository.
func (r *statusRepository) GetTimelinePostsByAccountId(
	ctx context.Context,
	accountID int,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return r.q.GetTimelinePostsByAccountId(ctx, int32(accountID))
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteByID(ctx context.Context, id int32) error {
	return r.q.DeleteStatusByID(ctx, id)
}
