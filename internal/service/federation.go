package service

import (
	"context"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/processing"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Object, error)
	ProcessIncoming(context.Context, *domain.ObjectOrLink) (func(context.Context) error, error)
}

type federationService struct {
	store     repo.Store
	processor processing.Processor
}

func NewFederationService(store repo.Store) FederationService {
	return &federationService{store: store}
}

var _ FederationService = (*federationService)(nil)

// GetLocalActor implements FederationService.
func (s *federationService) GetLocalActor(
	ctx context.Context,
	username string,
) (*domain.Object, error) {
	account, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	actor := ap.NewActor(nil)
	actor.SetObject(ap.Actor{
		ID:                account.Uri,
		Type:              "Person",
		PreferredUsername: account.Username,
		Inbox:             account.InboxUri,
		Outbox:            account.OutboxUri,
		Following:         account.FollowingUri,
		Followers:         account.FollowersUri,
	})
	return actor.GetRaw().Object, nil
}

// ProcessInbox implements FederationService.
func (s *federationService) ProcessIncoming(
	ctx context.Context,
	object *domain.ObjectOrLink,
) (func(context.Context) error, error) {
	activity := ap.NewActivity(object)
	if activity.GetValueType() != ap.ObjectType {
		return nil, errors.New("expected JSON object")
	}
	switch object.Object.Type {
	case "Create":
		return func(ctx context.Context) error {
			_, err := s.processor.Status(ap.NewNote(nil)).Get(ctx)
			return err
		}, nil
	case "Follow":
	case "Like":
	}
	return nil, nil
}
