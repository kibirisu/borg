package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type AccountRepository interface {
	Get(context.Context, db.GetAccountParams) (db.Account, error)
	GetByURI(context.Context, string) (db.Account, error)
	GetByID(context.Context, xid.ID) (db.Account, error)
	GetLocalByUsername(context.Context, string) (db.Account, error)
	Create(context.Context, db.CreateActorParams) (db.Account, error)
	GetFollowers(context.Context, xid.ID) ([]db.Account, error)
	GetFollowing(context.Context, xid.ID) ([]db.Account, error)
	GetPosts(context.Context, xid.ID) ([]db.GetStatusesByAccountIdRow, error)
	GetAccountRemoteFollowerInboxes(context.Context, xid.ID) ([]string, error)
}

type accountRepository struct {
	q *db.Queries
}

var _ AccountRepository = (*accountRepository)(nil)

// GetByURI implements AccountRepository.
func (r *accountRepository) GetByURI(ctx context.Context, uri string) (db.Account, error) {
	return r.q.GetActorByURI(ctx, uri)
}

// GetLocalByUsername implements AccountRepository.
func (r *accountRepository) GetLocalByUsername(
	ctx context.Context,
	username string,
) (db.Account, error) {
	return r.q.GetActor(ctx, username)
}

// Create implements AccountRepository.
func (r *accountRepository) Create(
	ctx context.Context,
	account db.CreateActorParams,
) (db.Account, error) {
	return r.q.CreateActor(ctx, account)
}

// Get implements AccountRepository.
func (r *accountRepository) Get(
	ctx context.Context,
	account db.GetAccountParams,
) (db.Account, error) {
	return r.q.GetAccount(ctx, account)
}

// GetById implements AccountRepository.
func (r *accountRepository) GetByID(
	ctx context.Context, id xid.ID,
) (db.Account, error) {
	return r.q.GetAccountById(ctx, id)
}

// GetFollowers implements AccountRepository.
func (r *accountRepository) GetFollowers(
	ctx context.Context, accountID xid.ID,
) ([]db.Account, error) {
	return r.q.GetAccountFollowers(ctx, accountID)
}

// GetFollowing implements AccountRepository.
func (r *accountRepository) GetFollowing(
	ctx context.Context, accountID xid.ID,
) ([]db.Account, error) {
	return r.q.GetAccountFollowing(ctx, accountID)
}

// GetPosts implements AccountRepository.
func (r *accountRepository) GetPosts(
	ctx context.Context,
	id xid.ID,
) ([]db.GetStatusesByAccountIdRow, error) {
	return r.q.GetStatusesByAccountId(ctx, id)
}

// GetAccountRemoteFollowerInboxes implements AccountRepository.
func (r *accountRepository) GetAccountRemoteFollowerInboxes(
	ctx context.Context,
	id xid.ID,
) ([]string, error) {
	return r.q.GetAccountRemoteFollowersInboxes(ctx, id)
}
