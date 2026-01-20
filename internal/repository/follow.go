package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRepository interface {
	AddFollowByActorURI(context.Context, db.AddFollowByActorURIParams) (string, error)
	GetLocalFollowByID(context.Context, xid.ID) (db.GetLocalFollowByIDRow, error)
	Create(context.Context, db.CreateFollowParams) (*db.Follow, error)
	CreateNew(context.Context, db.CreateFollowNewParams) error
	GetFollowerCollection(context.Context, string) (db.GetFollowerCollectionRow, error)
	GetFollowingCollection(context.Context, string) (db.GetFollowingCollectionRow, error)
	DeleteByID(context.Context, xid.ID) error
	GetByURI(context.Context, string) (db.Follow, error)
}

type followRepository struct {
	q *db.Queries
}

// AddFollowByActorURI implements FollowRepository.
func (r *followRepository) AddFollowByActorURI(
	ctx context.Context,
	follow db.AddFollowByActorURIParams,
) (string, error) {
	return r.q.AddFollowByActorURI(ctx, follow)
}

var _ FollowRepository = (*followRepository)(nil)

// GetLocalFollowByID implements FollowRepository.
func (r *followRepository) GetLocalFollowByID(
	ctx context.Context,
	id xid.ID,
) (db.GetLocalFollowByIDRow, error) {
	return r.q.GetLocalFollowByID(ctx, id)
}

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

// CreateNew implements FollowRepository.
func (r *followRepository) CreateNew(ctx context.Context, follow db.CreateFollowNewParams) error {
	return r.q.CreateFollowNew(ctx, follow)
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

// DeleteByID implements FollowRepository.
func (r *followRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteFollow(ctx, id)
}

func (r *followRepository) GetByURI(
	ctx context.Context, uri string,
) (db.Follow, error) {
	return r.q.GetFollowByURI(ctx, uri)
}
