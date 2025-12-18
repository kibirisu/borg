package service

import (
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	Foo()
}

type federationService struct {
	accounts repo.AccountRepository
	users    repo.UserRepository
}

func NewFederationService(q *db.Queries) FederationService {
	return &federationService{}
}

var _ FederationService = (*federationService)(nil)

// Foo implements FederationService.
func (s *federationService) Foo() {
	panic("unimplemented")
}
