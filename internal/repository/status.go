package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	Add(context.Context, db.AddStatusParams) error
	GetByID(context.Context, xid.ID) (db.Status, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetShares(context.Context, xid.ID) ([]db.Status, error)
	GetLocalStatuses(context.Context) ([]db.GetLocalStatusesRow, error)
	GetByIDWithMetadata(context.Context, xid.ID) (db.GetStatusByIdWithMetadataRow, error)
	GetSharedPostsByAccountId(context.Context, xid.ID) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountId(
		context.Context,
		xid.ID,
	) ([]db.GetTimelinePostsByAccountIdRow, error)
	DeleteByID(context.Context, xid.ID) error
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
func (r *statusRepository) GetByID(ctx context.Context, id xid.ID) (db.Status, error) {
	return r.q.GetStatusById(ctx, id)
}

// GetByURI implements StatusRepository.
func (r *statusRepository) GetByURI(ctx context.Context, uri string) (db.Status, error) {
	return r.q.GetStatusByURI(ctx, uri)
}

// GetById implements StatusRepository.
func (r *statusRepository) GetByIDWithMetadata(
	ctx context.Context,
	id xid.ID,
) (db.GetStatusByIdWithMetadataRow, error) {
	return r.q.GetStatusByIdWithMetadata(ctx, id)
}

// GetShares implements StatusRepository.
func (r *statusRepository) GetShares(ctx context.Context, id xid.ID) ([]db.Status, error) {
	return r.q.GetStatusShares(ctx, &id)
}

// GetLocalStatuses implements StatusRepository.
func (r *statusRepository) GetLocalStatuses(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return r.q.GetLocalStatuses(ctx)
}

// GetSharedPostsByAccountId implements StatusRepository.
func (r *statusRepository) GetSharedPostsByAccountId(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	return r.q.GetSharedPostsByAccountId(ctx, accountID)
}

// GetTimelinePostsByAccountId implements StatusRepository.
func (r *statusRepository) GetTimelinePostsByAccountId(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return r.q.GetTimelinePostsByAccountId(ctx, accountID)
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteStatusByID(ctx, id)
}
