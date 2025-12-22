package service

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Actor, error)
	CreateActor(context.Context, db.CreateActorParams) (*db.Account, error)
}

type federationService struct {
	store repo.Store
}

func NewFederationService(store repo.Store) FederationService {
	return &federationService{store}
}

var _ FederationService = (*federationService)(nil)

// GetLocalActor implements FederationService.
// not using anymore
func (s *federationService) GetLocalActor(
	ctx context.Context,
	username string,
) (*domain.Actor, error) {
	account, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	// we probably would implement mapper functions
	actor := domain.Actor{
		Context:           "https://www.w3.org/ns/activitystreams",
		ID:                account.Uri,
		Type:              "Person",
		PreferredUsername: account.Username,
		Inbox:             account.InboxUri,
		Outbox:            account.OutboxUri,
		Following:         account.FollowingUri,
		Followers:         account.FollowersUri,
	}
	return &actor, nil
}

// CreateActor implements FederationService.
func (s *federationService) CreateActor(
	ctx context.Context,
	actor db.CreateActorParams,
) (*db.Account, error) {
	account, err := s.store.Accounts().Create(ctx, actor)
	return &account, err
}
