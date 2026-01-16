package repository

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
)

type Store interface {
	Accounts() AccountRepository
	Users() UserRepository
	Follows() FollowRepository
	FollowRequests() FollowRequestRepository
	Statuses() StatusRepository
	Favourites() FavouriteRepository
	WithTX(context.Context, Tx) (any, error)
}

type Tx func(context.Context, Store) (any, error)

var _ Store = (*store)(nil)

type store struct {
	db *sql.DB
	q  *db.Queries
}

func New(ctx context.Context, url string) Store {
	db, q := db.GetDB(ctx, url)
	return &store{db, q}
}

// Accounts implements Store.
func (s *store) Accounts() AccountRepository {
	return &accountRepository{s.q}
}

// Users implements Store.
func (s *store) Users() UserRepository {
	return &userRepository{s.q}
}

func (s *store) Follows() FollowRepository {
	return &followRepository{s.q}
}

func (s *store) FollowRequests() FollowRequestRepository {
	return &followRequestRepository{s.q}
}

// Statuses implements Store.
func (s *store) Statuses() StatusRepository {
	return &statusRepository{s.q}
}

// Favourites implements Store.
func (s *store) Favourites() FavouriteRepository {
	return &favouriteRepository{s.q}
}

// WithTX implements Store.
func (s *store) WithTX(ctx context.Context, fn Tx) (any, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := s.q.WithTx(tx)
	store := store{q: qtx}

	res, err := fn(ctx, &store)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return res, err
}
