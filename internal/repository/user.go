package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type UserRepository interface {
	Create(context.Context, db.CreateUserParams) error
	GetByUsername(context.Context, string) (db.AuthDataRow, error)
}

type userRepository struct {
	q *db.Queries
}

var _ UserRepository = (*userRepository)(nil)

// Create implements UserRepository.
func (u *userRepository) Create(ctx context.Context, user db.CreateUserParams) error {
	return u.q.CreateUser(ctx, user)
}

// GetByUsername implements UserRepository.
func (u *userRepository) GetByUsername(
	ctx context.Context,
	username string,
) (db.AuthDataRow, error) {
	return u.q.AuthData(ctx, username)
}
