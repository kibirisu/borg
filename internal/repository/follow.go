package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRepository interface {
	Create(context.Context, db.CreateFollowParams) (*db.Follow, error)
}

type followRepository struct {
	q *db.Queries
}

var _ FollowRepository = (*followRepository)(nil)

// Create implements FollowRepository.
func (r *followRepository) Create(ctx context.Context, followCreate db.CreateFollowParams) (*db.Follow, error) {
	follow, err := r.q.CreateFollow(ctx, followCreate)
	if err != nil {
		return nil, err
	}
	return &follow, err
}
