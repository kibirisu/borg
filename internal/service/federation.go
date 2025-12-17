package service

import (
	"context"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	Foo(context.Context, *api.AuthForm)
}

type federationService struct {
	accounts repo.AccountRepository
	users    repo.UserRepository
}

func NewFederationService(q *db.Queries) FederationService {
	return &federationService{
		accounts: repo.NewAccountRepository(q),
		users:    repo.NewUserRepository(q),
	}
}

// Foo implements FederationService.
func (s *federationService) Foo(ctx context.Context, form *api.AuthForm) {
	_, _ = s.accounts.Create(ctx, db.CreateActorParams{})
}

var _ FederationService = (*federationService)(nil)
