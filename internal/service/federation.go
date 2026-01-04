package service

import (
	"context"
	"fmt"
	"log"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Actor, error)
	CreateActor(context.Context, db.CreateActorParams) (*db.Account, error)
	ProcessIncomingActivity(context.Context, *domain.Object) error
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
// func (s *federationService) GetLocalActor(
// 	ctx context.Context,
// 	username string,
// ) (*domain.Actor, error) {
// 	account, err := s.store.Accounts().GetLocalByUsername(ctx, username)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// we probably would implement mapper functions
// 	actor := domain.Actor{
// 		Context:           "https://www.w3.org/ns/activitystreams",
// 		ID:                account.Uri,
// 		Type:              "Person",
// 		PreferredUsername: account.Username,
// 		Inbox:             account.InboxUri,
// 		Outbox:            account.OutboxUri,
// 		Following:         account.FollowingUri,
// 		Followers:         account.FollowersUri,
// 	}
// 	return &actor, nil
// }

// CreateActor implements FederationService.
func (s *federationService) CreateActor(
	ctx context.Context,
	actor db.CreateActorParams,
) (*db.Account, error) {
	account, err := s.store.Accounts().Create(ctx, actor)
	return &account, err
}

// ProcessIncomingActivity implements FederationService.
func (s *federationService) ProcessIncomingActivity(ctx context.Context, raw *domain.Object) error {
	log.Printf("Processing incoming activity: %s (ID: %s)", raw.Type, raw.ID)

	switch raw.Type {
	case "Create":
		wrappedRaw := &domain.ObjectOrLink{Object: raw}
		activityWrapper := ap.NewCreateActivity(wrappedRaw)
		return s.handleCreate(ctx, activityWrapper.GetObject())

	case "Follow":
		return fmt.Errorf("Follow activity not implemented yet")

	case "Undo":
		return fmt.Errorf("Undo activity not implemented yet")

	default:
		return fmt.Errorf("unsupported activity type: %s", raw.Type)
	}
}

func (s *federationService) handleCreate(ctx context.Context, activity ap.Activity[ap.Note]) error {
	// TODO: Fetch remote actor if not exists (hit endpoint)
	actorID := activity.Actor.GetRaw().Object.ID

	// Extract Note content
	noteContent := activity.Object.GetObject().Content

	log.Printf("Authorized Actor %s created a Note: %s", actorID, noteContent)

	return nil
}
