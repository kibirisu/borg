package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRequestRepository interface {
	Create(context.Context, db.CreateFollowRequestParams) error
}

type followRequestRepository struct {
	q *db.Queries
}

var _ FollowRequestRepository = (*followRequestRepository)(nil)

// Create implements FollowRepository.
func (r *followRequestRepository) Create(
	ctx context.Context,
	request db.CreateFollowRequestParams,
) error {
	return r.q.CreateFollowRequest(ctx, request)
}
