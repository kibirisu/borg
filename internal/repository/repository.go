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
	Statuses() StatusRepository
	Favourites() FavouriteRepository
}

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

// Statuses implements Store.
func (s *store) Statuses() StatusRepository {
	return &statusRepository{s.q}
}

// Favourites implements Store.
func (s *store) Favourites() FavouriteRepository {
	return &favouriteRepository{s.q}
}
