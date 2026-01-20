package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	AddByActorURI(context.Context, db.AddStatusByActorURIParams) error
	GetByID(context.Context, db.GetStatusByIDParams) (db.GetStatusByIDRow, error)
	GetByAccountID(
		context.Context,
		db.GetStatusesByAccountIDParams,
	) ([]db.GetStatusesByAccountIDRow, error)
	CreateNew(context.Context, db.CreateStatusNewParams) (db.Status, error)
	ReblogStatus(context.Context, db.CreateReblogParams) (db.Status, error)
	DeleteByIDNew(context.Context, xid.ID) (db.Status, error)
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	GetReplies(context.Context, xid.ID, xid.ID) ([]db.GetStatusRepliesRow, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetLocalByID(context.Context, xid.ID) (db.Status, error)
	GetHomeTimelineByAccountID(context.Context, xid.ID) ([]db.GetTimelinePostsByAccountIdRow, error)
	DeleteByID(context.Context, xid.ID) error
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// AddByActorURI implements StatusRepository.
func (r *statusRepository) AddByActorURI(
	ctx context.Context,
	status db.AddStatusByActorURIParams,
) error {
	return r.q.AddStatusByActorURI(ctx, status)
}

// GetByID implements StatusRepository.
func (r *statusRepository) GetByID(
	ctx context.Context,
	ids db.GetStatusByIDParams,
) (db.GetStatusByIDRow, error) {
	return r.q.GetStatusByID(ctx, ids)
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

// DeleteByIDNew implements StatusRepository.
func (r *statusRepository) DeleteByIDNew(ctx context.Context, id xid.ID) (db.Status, error) {
	return r.q.DeleteStatusByIDNew(ctx, id)
}

// Create implements StatusRepository.
func (r *statusRepository) Create(
	ctx context.Context,
	status db.CreateStatusParams,
) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}

// GetReplies implements StatusRepository.
func (r *statusRepository) GetReplies(
	ctx context.Context,
	accountID xid.ID,
	statusID xid.ID,
) ([]db.GetStatusRepliesRow, error) {
	param := db.GetStatusRepliesParams{
		InReplyToID: &statusID,
		AccountID:   accountID,
	}
	return r.q.GetStatusReplies(ctx, param)
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

// GetLocalByID implements StatusRepository.
func (r *statusRepository) GetLocalByID(ctx context.Context, id xid.ID) (db.Status, error) {
	return r.q.GetLocalStatusByID(ctx, id)
}

// GetHomeTimelineByAccountID implements StatusRepository.
func (r *statusRepository) GetHomeTimelineByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return r.q.GetTimelinePostsByAccountId(ctx, id)
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteStatusByID(ctx, id)
}
