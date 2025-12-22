package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRepository interface {
	Create(context.Context, db.CreateFollowParams) error
}

type followRepository struct {
	q *db.Queries
}

var _ FollowRepository = (*followRepository)(nil)

// Create implements FollowRepository.
func (r *followRepository) Create(ctx context.Context, follow db.CreateFollowParams) error {
	return r.q.CreateFollow(ctx, follow)
}
