package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type AccountRepository interface {
	Get(context.Context, db.GetAccountParams) (db.Account, error)
	GetLocalByUsername(context.Context, string) (db.Account, error)
	Create(context.Context, db.CreateActorParams) (db.Account, error)
}

type accountRepository struct {
	q *db.Queries
}

var _ AccountRepository = (*accountRepository)(nil)

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
