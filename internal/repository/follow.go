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
func (u *followRepository) Create(ctx context.Context, follow db.CreateFollowParams) error {
	return u.q.CreateFollow(ctx, follow)
}
