package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type AccountRepository interface {
	GetByURI(context.Context, string) (db.Account, error)
	GetByID(context.Context, xid.ID) (db.GetAccountByIDRow, error)
	GetFollowersByAccountID(context.Context, xid.ID) ([]db.GetFollowersByAccountIDRow, error)
	GetFollowingByAccountID(context.Context, xid.ID) ([]db.GetFollowingByAccountIDRow, error)
	GetLocalByUsername(context.Context, string) (db.Account, error)
	Create(context.Context, db.CreateActorParams) (db.Account, error)
	GetFollowers(context.Context, xid.ID) ([]db.Account, error)
	GetFollowing(context.Context, xid.ID) ([]db.Account, error)
	GetAccountRemoteFollowerInboxes(context.Context, xid.ID) ([]string, error)
	GetAccountInbox(context.Context, xid.ID) (string, error)
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

// GetById implements AccountRepository.
func (r *accountRepository) GetByID(
	ctx context.Context, id xid.ID,
) (db.GetAccountByIDRow, error) {
	return r.q.GetAccountByID(ctx, id)
}

// GetFollowersByAccountID implements AccountRepository.
func (r *accountRepository) GetFollowersByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetFollowersByAccountIDRow, error) {
	return r.q.GetFollowersByAccountID(ctx, id)
}

// GetFollowingByAccountID implements AccountRepository.
func (r *accountRepository) GetFollowingByAccountID(
	ctx context.Context,
	id xid.ID,
) ([]db.GetFollowingByAccountIDRow, error) {
	return r.q.GetFollowingByAccountID(ctx, id)
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

// GetAccountRemoteFollowerInboxes implements AccountRepository.
func (r *accountRepository) GetAccountRemoteFollowerInboxes(
	ctx context.Context,
	id xid.ID,
) ([]string, error) {
	return r.q.GetAccountRemoteFollowersInboxes(ctx, id)
}

// GetAccountInbox implements AccountRepository.
func (r *accountRepository) GetAccountInbox(ctx context.Context, id xid.ID) (string, error) {
	return r.q.GetAccountInbox(ctx, id)
}
