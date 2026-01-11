package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRepository interface {
	Create(context.Context, db.CreateFollowParams) (*db.Follow, error)
	GetFollowerCollection(context.Context, string) (db.GetFollowerCollectionRow, error)
	GetFollowingCollection(context.Context, string) (db.GetFollowingCollectionRow, error)
	GetByURI(context.Context, string) (db.Follow, error)
}

type followRepository struct {
	q *db.Queries
}

var _ FollowRepository = (*followRepository)(nil)

// Create implements FollowRepository.
func (r *followRepository) Create(
	ctx context.Context,
	followCreate db.CreateFollowParams,
) (*db.Follow, error) {
	follow, err := r.q.CreateFollow(ctx, followCreate)
	if err != nil {
		return nil, err
	}
	return &follow, err
}

// GetFollowerCollection implements FollowRepository.
func (r *followRepository) GetFollowerCollection(
	ctx context.Context,
	username string,
) (db.GetFollowerCollectionRow, error) {
	return r.q.GetFollowerCollection(ctx, username)
}

// GetFollowingCollection implements FollowRepository.
func (r *followRepository) GetFollowingCollection(
	ctx context.Context,
	username string,
) (db.GetFollowingCollectionRow, error) {
	return r.q.GetFollowingCollection(ctx, username)
}

func (r *followRepository) GetByURI(
	ctx context.Context, uri string,
) (db.Follow, error) {
	return r.q.GetFollowByURI(ctx, uri)
}
