package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	GetByIDNew(context.Context, db.GetStatusByIDNewParams) (db.GetStatusByIDNewRow, error)
	GetByAccountID(
		context.Context,
		db.GetStatusesByAccountIDParams,
	) ([]db.GetStatusesByAccountIDRow, error)
	CreateNew(context.Context, db.CreateStatusNewParams) (db.Status, error)
	ReblogStatus(context.Context, db.CreateReblogParams) (db.Status, error)
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	GetByID(context.Context, xid.ID) (db.Status, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetByIDWithMetadata(context.Context, xid.ID) (db.GetStatusByIdWithMetadataRow, error)
	GetSharedPostsByAccountID(context.Context, xid.ID) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountID(
		context.Context,
		xid.ID,
	) ([]db.GetTimelinePostsByAccountIdRow, error)
	DeleteByID(context.Context, xid.ID) error
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// GetByIDNew implements StatusRepository.
func (r *statusRepository) GetByIDNew(
	ctx context.Context,
	ids db.GetStatusByIDNewParams,
) (db.GetStatusByIDNewRow, error) {
	return r.q.GetStatusByIDNew(ctx, ids)
}

// CreateNew implements StatusRepository.
func (r *statusRepository) CreateNew(
	ctx context.Context,
	status db.CreateStatusNewParams,
) (db.Status, error) {
	return r.q.CreateStatusNew(ctx, status)
}

// ReblogStatus implements StatusRepository.
func (r *statusRepository) ReblogStatus(
	ctx context.Context,
	reblog db.CreateReblogParams,
) (db.Status, error) {
	return r.q.CreateReblog(ctx, reblog)
}

// Create implements StatusRepository.
func (r *statusRepository) Create(
	ctx context.Context,
	status db.CreateStatusParams,
) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}

// GetById implements StatusRepository.
func (r *statusRepository) GetByID(ctx context.Context, id xid.ID) (db.Status, error) {
	return r.q.GetStatusById(ctx, id)
}

// GetByAccountID implements StatusRepository.
func (r *statusRepository) GetByAccountID(
	ctx context.Context,
	ids db.GetStatusesByAccountIDParams,
) ([]db.GetStatusesByAccountIDRow, error) {
	return r.q.GetStatusesByAccountID(ctx, ids)
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

// GetSharedPostsByAccountId implements StatusRepository.
func (r *statusRepository) GetSharedPostsByAccountID(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	return r.q.GetSharedPostsByAccountId(ctx, accountID)
}

// GetTimelinePostsByAccountId implements StatusRepository.
func (r *statusRepository) GetTimelinePostsByAccountID(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return r.q.GetTimelinePostsByAccountId(ctx, accountID)
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteStatusByID(ctx, id)
}
