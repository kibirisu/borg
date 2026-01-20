package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type FollowRepository interface {
	Create(context.Context, db.CreateFollowNewParams) error
	AddByRequestURI(context.Context, string) error
	AddByActorURI(context.Context, db.AddFollowByActorURIParams) (string, error)
	GetLocalFollowByID(context.Context, xid.ID) (db.GetLocalFollowByIDRow, error)
	GetFollowerCollection(context.Context, string) (db.GetFollowerCollectionRow, error)
	GetFollowingCollection(context.Context, string) (db.GetFollowingCollectionRow, error)
	DeleteByID(context.Context, xid.ID) error
	DeleteByURI(context.Context, string) error
}

type followRepository struct {
	q *db.Queries
}

var _ FollowRepository = (*followRepository)(nil)

// AddByActorURI implements FollowRepository.
func (r *followRepository) AddByActorURI(
	ctx context.Context,
	follow db.AddFollowByActorURIParams,
) (string, error) {
	return r.q.AddFollowByActorURI(ctx, follow)
}

// AddByRequestURI implements FollowRepository.
func (r *followRepository) AddByRequestURI(ctx context.Context, uri string) error {
	return r.q.AddFollowByRequestURI(ctx, uri)
}

// GetLocalFollowByID implements FollowRepository.
func (r *followRepository) GetLocalFollowByID(
	ctx context.Context,
	id xid.ID,
) (db.GetLocalFollowByIDRow, error) {
	return r.q.GetLocalFollowByID(ctx, id)
}

// Create implements FollowRepository.
func (r *followRepository) Create(ctx context.Context, follow db.CreateFollowNewParams) error {
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
	return r.q.DeleteFollowByID(ctx, id)
}

// DeleteByURI implements FollowRepository.
func (r *followRepository) DeleteByURI(ctx context.Context, uri string) error {
	return r.q.DeleteFollowByURI(ctx, uri)
}
