package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Object, error)
	CreateActor(context.Context, db.CreateActorParams) (*db.Account, error)
	ProcessIncoming(context.Context, *domain.ObjectOrLink) error
}

type federationService struct {
	store  repo.Store
	client transport.Client
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

// CreateActor implements FederationService.
func (s *federationService) CreateActor(
	ctx context.Context,
	actor db.CreateActorParams,
) (*db.Account, error) {
	account, err := s.store.Accounts().Create(ctx, actor)
	return &account, err
}

// ProcessInbox implements FederationService.
func (s *federationService) ProcessIncoming(
	ctx context.Context,
	object *domain.ObjectOrLink,
) error {
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

	attributedTo := getActorURI(note.AttributedTo)
	account, err := s.store.Accounts().GetByURI(ctx, attributedTo)
	if err != nil {
		obj, err := s.client.Get(attributedTo)
		if err != nil {
			return err
		}
		actor := ap.NewActor(obj)
		actorData := actor.GetObject()
		account, err = s.store.Accounts().Create(ctx, db.CreateActorParams{
			Username:     actorData.PreferredUsername,
			Uri:          actorData.ID,
			DisplayName:  sql.NullString{},
			Domain:       sql.NullString{},
			InboxUri:     actorData.Inbox,
			OutboxUri:    actorData.Outbox,
			Url:          "nope",
			FollowersUri: actorData.Followers,
			FollowingUri: actorData.Following,
		})
		if err != nil {
			return err
		}
	}
	inReplyTo := sql.NullInt32{}
	var inReplyToURI string
	switch note.InReplyTo.GetValueType() {
	case ap.LinkType:
		inReplyToURI = note.InReplyTo.GetURI()
	case ap.NullType:
	case ap.ObjectType:
		inReplyToURI = note.InReplyTo.GetObject().ID
	default:
		panic("unexpected ap.ValueType")
	}

	if inReplyToURI != "" {
		parentStatus, err := s.store.Statuses().GetByURI(ctx, inReplyToURI)
		if err != nil {
			obj, err := s.client.Get(inReplyToURI)
			if err != nil {
				return err
			}
			status := ap.NewNote(obj)
			statusData := status.GetObject()
			parentStatus, err = s.store.Statuses().Create(ctx, db.CreateStatusParams{
				Uri:         statusData.ID,
				Url:         "nope",
				Local:       sql.NullBool{},
				Content:     statusData.Content,
				AccountID:   0, // TODO: check again for actor
				InReplyToID: inReplyTo,
				ReblogOfID:  sql.NullInt32{},
			})
			if err != nil {
				return err
			}
		}
		inReplyTo = sql.NullInt32{
			Int32: parentStatus.ID,
			Valid: true,
		}
	}

	return s.store.Statuses().Add(ctx, db.AddStatusParams{
		Uri:         note.ID,
		Url:         "nope",
		Content:     note.Content,
		AccountID:   account.ID,
		InReplyToID: inReplyTo,
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
