package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type StatusRepository interface {
	AddReblog(context.Context, db.AddReblogParams) (db.Status, error)
	GetByID(context.Context, db.GetStatusByIDParams) (db.GetStatusByIDRow, error)
	GetByAccountID(
		context.Context,
		db.GetStatusesByAccountIDParams,
	) ([]db.GetStatusesByAccountIDRow, error)
	Create(context.Context, db.CreateStatusParams) (db.Status, error)
	CreateReply(context.Context, db.CreateReplyParams) (db.CreateReplyRow, error)
	ReblogStatus(context.Context, db.CreateReblogParams) (db.CreateReblogRow, error)
	Add(context.Context, db.AddStatusParams) (db.Status, error)
	AddReply(context.Context, db.AddReplyParams) (db.Status, error)
	GetReplies(context.Context, xid.ID, xid.ID) ([]db.GetStatusRepliesRow, error)
	GetByURI(context.Context, string) (db.Status, error)
	GetLocalByID(context.Context, xid.ID, xid.ID) (db.GetLocalStatusByIDRow, error)
	GetHomeTimelineByAccountID(context.Context, xid.ID) ([]db.GetTimelinePostsByAccountIdRow, error)
	GetFavouriteByAccountID(context.Context, xid.ID) ([]db.GetFavouritePostsByAccountIdRow, error)
	GetRebloggedByAccountID(context.Context, xid.ID) ([]db.GetRebloggedPostsByAccountIdRow, error)
	DeleteAnnounceByID(context.Context, xid.ID) error
	DeleteReblogByStatusID(
		context.Context,
		xid.ID,
		xid.ID,
	) (db.DeleteReblogByStatusIDRow, error)
}

type statusRepository struct {
	q *db.Queries
}

var _ StatusRepository = (*statusRepository)(nil)

// AddReblog implements StatusRepository.
func (r *statusRepository) AddReblog(
	ctx context.Context,
	reblog db.AddReblogParams,
) (db.Status, error) {
	return r.q.AddReblog(ctx, reblog)
}

// GetByID implements StatusRepository.
func (r *statusRepository) GetByID(
	ctx context.Context,
	ids db.GetStatusByIDParams,
) (db.GetStatusByIDRow, error) {
	return r.q.GetStatusByID(ctx, ids)
}

// CreateReply implements StatusRepository.
func (r *statusRepository) CreateReply(
	ctx context.Context,
	reply db.CreateReplyParams,
) (db.CreateReplyRow, error) {
	return r.q.CreateReply(ctx, reply)
}

// Create implements StatusRepository.
func (r *statusRepository) Create(
	ctx context.Context,
	status db.CreateStatusParams,
) (db.Status, error) {
	return r.q.CreateStatus(ctx, status)
}

// ReblogStatus implements StatusRepository.
func (r *statusRepository) ReblogStatus(
	ctx context.Context,
	reblog db.CreateReblogParams,
) (db.CreateReblogRow, error) {
	return r.q.CreateReblog(ctx, reblog)
}

// Add implements StatusRepository.
func (r *statusRepository) Add(
	ctx context.Context,
	status db.AddStatusParams,
) (db.Status, error) {
	return r.q.AddStatus(ctx, status)
}

// AddReply implements StatusRepository.
func (r *statusRepository) AddReply(
	ctx context.Context,
	reply db.AddReplyParams,
) (db.Status, error) {
	return r.q.AddReply(ctx, reply)
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
func (r *statusRepository) GetLocalByID(
	ctx context.Context,
	statusID, accountID xid.ID,
) (db.GetLocalStatusByIDRow, error) {
	return r.q.GetLocalStatusByID(ctx, db.GetLocalStatusByIDParams{
		StatusID:  statusID,
		AccountID: accountID,
	})
}

// GetHomeTimelineByAccountID implements StatusRepository.
func (r *statusRepository) GetHomeTimelineByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return r.q.GetTimelinePostsByAccountId(ctx, id)
}

// GetFavouriteByAccountID implements StatusRepository.
func (r *statusRepository) GetFavouriteByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetFavouritePostsByAccountIdRow, error) {
	return r.q.GetFavouritePostsByAccountId(ctx, id)
}

// GetRebloggedByAccountID implements StatusRepository.
func (r *statusRepository) GetRebloggedByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetRebloggedPostsByAccountIdRow, error) {
	return r.q.GetRebloggedPostsByAccountId(ctx, id)
}

// DeleteByURI implements StatusRepository.
func (r *statusRepository) DeleteAnnounceByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteAnnounceByID(ctx, id)
}

// DeleteReblogByStatusID implements StatusRepository.
func (r *statusRepository) DeleteReblogByStatusID(
	ctx context.Context,
	statusID xid.ID,
	accountID xid.ID,
) (db.DeleteReblogByStatusIDRow, error) {
	return r.q.DeleteReblogByStatusID(ctx, db.DeleteReblogByStatusIDParams{
		AccountID:  accountID,
		ReblogOfID: &statusID,
	})
}
