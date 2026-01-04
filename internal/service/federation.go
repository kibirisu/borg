package service

import (
	"context"
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
	ProcessInbox(context.Context, *domain.ObjectOrLink) error
}

type federationService struct {
	store repo.Store
}

func NewFederationService(store repo.Store) FederationService {
	return &federationService{store}
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

// CreateActor implements FederationService.
func (s *federationService) CreateActor(
	ctx context.Context,
	actor db.CreateActorParams,
) (*db.Account, error) {
	account, err := s.store.Accounts().Create(ctx, actor)
	return &account, err
}

// ProcessInbox implements FederationService.
func (s *federationService) ProcessInbox(ctx context.Context, object *domain.ObjectOrLink) error {
	activity := ap.NewActivity(object)
	if activity.GetValueType() != ap.ObjectType {
		return errors.New("expected JSON object")
	}
	activityData := activity.GetObject()
	switch activityData.Type {
	case "Create":
		return s.addNote(ctx, activity)
	case "Follow":
		_ = s.AddFollow(ctx, activity)
	case "Like":
	}
	return nil
}

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

func (s *federationService) addNote(ctx context.Context, activity ap.Activiter[any]) error {
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

	attributedTo := getActorURI(note.AttributedTo)

	err := s.store.Statuses().AddFrom(ctx, db.AddStatusFromParams{
		Uri:         note.ID,
		Url:         "TODO",
		Content:     note.Content,
		Uri_2:       attributedTo,
		InReplyToID: sql.NullInt32{},
		ReblogOfID:  sql.NullInt32{},
	})
	if err != nil {
		return err
	}

	return s.store.Statuses().Add(ctx, db.AddStatusParams{
		Uri:         note.ID,
		Url:         "TODO",
		Content:     note.Content,
		AccountID:   0, // TODO
		InReplyToID: sql.NullInt32{},
		ReblogOfID:  sql.NullInt32{},
	})
}

func getActorURI(object ap.Actorer) string {
	switch object.GetValueType() {
	case ap.LinkType:
		return object.GetURI()
	case ap.ObjectType:
		return object.GetObject().ID
	case ap.NullType:
		fallthrough
	default:
		panic("unexpected ap.ValueType")
	}
}

func (s *federationService) getActor() {
}
