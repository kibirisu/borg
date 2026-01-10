package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/domain"
	proc "github.com/kibirisu/borg/internal/processing"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/worker"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Object, error)
	GetActorFollowers(context.Context, string, *int) (*domain.Object, error)
	GetActorFollowing(context.Context, string, *int) (*domain.Object, error)
	ProcessIncoming(context.Context, *domain.ObjectOrLink) (worker.Job, error)
}

type federationService struct {
	store     repo.Store
	processor proc.Processor
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

// GetActorFollowers implements FederationService.
func (s *federationService) GetActorFollowers(
	ctx context.Context,
	username string,
	pageSelection *int,
) (*domain.Object, error) {
	data, err := s.store.Follows().GetFollowerCollection(ctx, username)
	if err != nil {
		return nil, err
	}
	local, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if pageSelection != nil {
		followers, _ := s.store.Accounts().GetFollowers(ctx, int(local.ID))
		links := make([]ap.Objecter[ap.Actor], 0, len(followers))
		for _, acc := range followers {
			actor := ap.NewActor(nil)
			actor.SetLink(acc.Uri)
			links = append(links, actor)
			log.Println(acc.Uri)
		}
		page := ap.NewActorCollectionPage(nil)
		page.SetObject(ap.CollectionPage[ap.Actor]{
			ID:     fmt.Sprintf("%s?page=%d", data.FollowersUri, *pageSelection),
			Type:   "OrderedCollectionPage",
			PartOf: ap.NewActorCollection(&domain.ObjectOrLink{Link: &data.FollowersUri}),
			Items:  links,
			Next:   ap.NewActorCollectionPage(nil),
		})

		return page.GetRaw().Object, nil
	}
	collection := ap.NewActorCollection(nil)
	page := ap.NewActorCollectionPage(nil)
	page.SetLink(fmt.Sprintf("%s?page=1", data.FollowersUri))
	collection.SetObject(ap.Collection[ap.Actor]{
		ID:    data.FollowersUri,
		Type:  "OrderedCollection",
		First: page,
	})
	return collection.GetRaw().Object, nil
}

// GetActorFollowing implements FederationService.
func (s *federationService) GetActorFollowing(
	ctx context.Context,
	username string,
	pageSelection *int,
) (*domain.Object, error) {
	data, err := s.store.Follows().GetFollowingCollection(ctx, username)
	if err != nil {
		return nil, err
	}
	local, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if pageSelection != nil {
		following, err := s.store.Accounts().GetFollowing(ctx, int(local.ID))
		if err != nil {
			return nil, err
		}
		links := make([]ap.Objecter[ap.Actor], len(following))
		for idx, acc := range following {
			actor := ap.NewActor(nil)
			actor.SetLink(acc.Uri)
			links[idx] = actor
			log.Println(acc.Uri)
		}
		page := ap.NewActorCollectionPage(nil)
		page.SetObject(ap.CollectionPage[ap.Actor]{
			ID:     fmt.Sprintf("%s?page=%d", data.FollowingUri, *pageSelection),
			Type:   "OrderedCollectionPage",
			PartOf: ap.NewActorCollection(&domain.ObjectOrLink{Link: &data.FollowingUri}),
			Items:  links,
			Next:   ap.NewActorCollectionPage(nil),
		})

		return page.GetRaw().Object, nil
	}
	collection := ap.NewActorCollection(nil)
	page := ap.NewActorCollectionPage(nil)
	page.SetLink(fmt.Sprintf("%s?page=1", data.FollowingUri))
	collection.SetObject(ap.Collection[ap.Actor]{
		ID:    data.FollowingUri,
		Type:  "OrderedCollection",
		First: page,
	})
	return collection.GetRaw().Object, nil
}

// ProcessInbox implements FederationService.
func (s *federationService) ProcessIncoming(
	ctx context.Context,
	object *domain.ObjectOrLink,
) (worker.Job, error) {
	if object.GetType() != domain.ObjectType {
		return nil, errors.New("expected JSON object")
	}
	switch object.Object.Type {
	case "Create":
		return func(ctx context.Context) error {
			_, err := s.processor.LookupStatus(ctx, ap.NewNote(object.Object.ActivityObject))
			return err
		}, nil
	case "Follow":
		return func(ctx context.Context) error {
			return s.processor.AcceptFollow(ctx, ap.NewFollowActivity(object))
		}, nil
	case "Accept":
		fallthrough
	case "Undo":
		fallthrough
	default:
		return nil, errors.New("unsupported Activity type")
	}
}
