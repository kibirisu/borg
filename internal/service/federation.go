package service

import (
	"context"
	"fmt"
	"log"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	repo "github.com/kibirisu/borg/internal/repository"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Object, error)
	CreateActor(context.Context, db.CreateActorParams) (*db.Account, error)
	ProcessIncomingActivity(context.Context, *domain.Object) error
	ProcessInbox(context.Context, *domain.ObjectOrLink) error
	AddFollow(context.Context, ap.Activiter[any]) error
	AddNote(context.Context, ap.Activiter[any]) error
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
// ProcessInbox implements FederationService.
func (s *federationService) ProcessInbox(ctx context.Context, object *domain.ObjectOrLink) error {
	activity := ap.NewActivity(object)
	if activity.GetValueType() != ap.ObjectType {
		return errors.New("expected JSON object")
	}
	activityData := activity.GetObject()
	switch activityData.Type {
	case "Create":
		return s.AddNote(ctx, activity)
	case "Follow":
		_ = s.AddFollow(ctx, activity)
	}
	return nil
}

// AddFollow implements FederationService.
func (s *federationService) AddFollow(ctx context.Context, activity ap.Activiter[any]) error {
	objectData := activity.GetObject().Object.GetRaw()
	object := ap.NewActor(objectData)
	var actorURI string
	switch object.GetValueType() {
	case ap.LinkType:
		actorURI = object.GetURI()
	case ap.ObjectType:

	case ap.NullType:
	default:
		panic("unexpected ap.ValueType")
	}
	_ = actorURI
	panic("unimplemented")
}

// CreateNote implements FederationService.
func (s *federationService) AddNote(ctx context.Context, activity ap.Activiter[any]) error {
	objectData := activity.GetObject().Object.GetRaw()
	object := ap.NewNote(objectData)
	if object.GetValueType() != ap.ObjectType {
		return errors.New("expected JSON object")
	}
	note := object.GetObject()
	if note.Type != "Note" {
		return errors.New("expected Note object")
	}

	// We may make db rows searchable by AP object URI for smoothness
	// Since we receiving that note, the account shall be present in db
	// Otherwise, it may be DM, so we fetch remote actor

	return s.store.Statuses().Add(ctx, db.AddStatusParams{
		Uri:         note.ID,
		Url:         "TODO",
		Content:     note.Content,
		AccountID:   0, // TODO
		InReplyToID: sql.NullInt32{},
		ReblogOfID:  sql.NullInt32{},
	})
}
