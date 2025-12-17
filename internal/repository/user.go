package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type UserRepository interface {
	Create(context.Context, db.CreateUserParams) error
}

type userRepository struct {
	q *db.Queries
}

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q}
}

// Create implements UserRepository.
func (u *userRepository) Create(ctx context.Context, user db.CreateUserParams) error {
	return u.q.CreateUser(ctx, user)
}
