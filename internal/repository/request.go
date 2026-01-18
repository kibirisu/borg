package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRequestRepository interface {
	Create(context.Context, db.CreateFollowRequestParams) (db.CreateFollowRequestRow, error)
	DeleteByTargetAccountID(
		context.Context,
		db.DeleteFollowRequestByAccountIDParams,
	) (db.DeleteFollowRequestByAccountIDRow, error)
}

type followRequestRepository struct {
	q *db.Queries
}

var _ FollowRequestRepository = (*followRequestRepository)(nil)

// Create implements FollowRepository.
func (r *followRequestRepository) Create(
	ctx context.Context,
	request db.CreateFollowRequestParams,
) (db.CreateFollowRequestRow, error) {
	return r.q.CreateFollowRequest(ctx, request)
}

// DeleteByTargetAccountID implements FollowRequestRepository.
func (r *followRequestRepository) DeleteByTargetAccountID(
	ctx context.Context,
	ids db.DeleteFollowRequestByAccountIDParams,
) (db.DeleteFollowRequestByAccountIDRow, error) {
	return r.q.DeleteFollowRequestByAccountID(ctx, ids)
}
