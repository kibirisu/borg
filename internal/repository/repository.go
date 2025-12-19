package repository

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
)

type Store interface {
	Accounts() AccountRepository
	Users() UserRepository
}

var _ Store = (*store)(nil)

type store struct {
	db *sql.DB
	q  *db.Queries
}

func NewStore(url string) Store {
	db, q := db.GetDB(context.Background(), url)
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
