package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type AccountRepository interface {
	GetLocalByUsername(context.Context, string) (db.Account, error)
	Create(context.Context, db.CreateActorParams) (db.Account, error)
}

type accountRepository struct {
	q *db.Queries
}

var _ AccountRepository = (*accountRepository)(nil)

func NewAccountRepository(q *db.Queries) AccountRepository {
	return &accountRepository{q}
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
