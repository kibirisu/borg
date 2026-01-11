package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/domain"
	proc "github.com/kibirisu/borg/internal/processing"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/util"
	"github.com/kibirisu/borg/internal/worker"
)

type FederationService interface {
	LookupRemoteActor(context.Context, *util.HandleInfo) (*domain.Actor, error)
	GetLocalActor(context.Context, string) (*domain.Object, error)
	GetActorFollowers(context.Context, string, *int) (*domain.Object, error)
	GetActorFollowing(context.Context, string, *int) (*domain.Object, error)
	GetStatus(context.Context, string) (*domain.Object, error)
	GetLike(context.Context, string) (*domain.Object, error)
	GetFollow(context.Context, string) (*domain.Object, error)
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

// GetStatus implements FederationService.
func (s *federationService) GetStatus(
	ctx context.Context,
	uri string,
) (*domain.Object, error) {
	statusDB, err := s.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		return nil, err
	}
	reply := ap.NewNote(nil)
	if statusDB.InReplyToID.Valid {
		inreply, err := s.store.Statuses().GetByID(ctx, int(statusDB.InReplyToID.Int32))
		if err != nil {
			return nil, err
		}

		reply.SetLink(inreply.Uri)
	}
	accountDB, err := s.store.Accounts().GetByID(ctx, int(statusDB.AccountID))
	if err != nil {
		return nil, err
	}
	actor := ap.NewActor(nil)
	actor.SetLink(accountDB.Uri)

	note := ap.NewNote(nil)
	collection := ap.NewNoteCollection(nil)
	note.SetObject(ap.Note{
		ID:           statusDB.Uri,
		Type:         "Note",
		Content:      statusDB.Content,
		InReplyTo:    reply,
		Published:    statusDB.CreatedAt,
		AttributedTo: actor,
		To:           []string{},
		CC:           []string{accountDB.FollowersUri},
		Replies:      collection,
	})

	return note.GetRaw().Object, nil
}

// GetLike implements FederationService.
func (s *federationService) GetLike(
	ctx context.Context,
	uri string,
) (*domain.Object, error) {
	likeDB, err := s.store.Favourites().GetByURI(ctx, uri)
	if err != nil {
		return nil, err
	}
	accountDB, err := s.store.Accounts().GetByID(ctx, int(likeDB.AccountID))
	if err != nil {
		return nil, err
	}
	postDB, err := s.store.Statuses().GetByID(ctx, int(likeDB.StatusID))
	if err != nil {
		return nil, err
	}
	actor := ap.NewActor(nil)
	actor.SetLink(accountDB.Uri)

	note := ap.NewNote(nil)
	note.SetLink(postDB.Uri)

	like := ap.NewLikeActivity(nil)
	like.SetObject(ap.Activity[ap.Note]{
		ID:     likeDB.Uri,
		Type:   "Like",
		Actor:  actor,
		Object: note,
	})

	return like.GetRaw().Object, nil
}

// GetFollow implements FederationService.
func (s *federationService) GetFollow(
	ctx context.Context,
	uri string,
) (*domain.Object, error) {
	followDB, err := s.store.Follows().GetByURI(ctx, uri)
	if err != nil {
		return nil, err
	}
	followerDB, err := s.store.Accounts().GetByID(ctx, int(followDB.AccountID))
	if err != nil {
		return nil, err
	}
	followeeDB, err := s.store.Accounts().GetByID(ctx, int(followDB.TargetAccountID))
	if err != nil {
		return nil, err
	}
	actorFollowed := ap.NewActor(nil)
	actorFollowed.SetLink(followeeDB.Uri)

	actorFollowing := ap.NewActor(nil)
	actorFollowing.SetLink(followerDB.Uri)

	follow := ap.NewFollowActivity(nil)
	follow.SetObject(ap.Activity[ap.Actor]{
		ID:     followDB.Uri,
		Type:   "Follow",
		Actor:  actorFollowed,
		Object: actorFollowing,
	})

	return follow.GetRaw().Object, nil
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
	case "Announce":
		return func(ctx context.Context) error {
			_, err := s.processor.AnnounceStatus(ctx, ap.NewAnnounceActivity(object))
			return err
		}, nil
	case "Like":
		return func(ctx context.Context) error {
			_, err := s.processor.LikeStatus(ctx, ap.NewLikeActivity(object))
			return err
		}, nil
	case "Accept":
		fallthrough
	case "Undo":
		return s.processUndo(object)
	case "Delete":
		return s.processDelete(object)
	default:
		return nil, errors.New("unsupported Activity type")
	}
}

func (s *federationService) processUndo(object *domain.ObjectOrLink) (worker.Job, error) {
	switch object.Object.Type {
	case "Announce":
		return func(ctx context.Context) error {
			status, err := s.processor.AnnounceStatus(
				ctx,
				ap.NewAnnounceActivity(object.Object.ActivityObject),
			)
			if err != nil {
				return err
			}
			return s.store.Statuses().DeleteByID(ctx, status.ID)
		}, nil
	case "Like":
		return func(ctx context.Context) error {
			favourite, err := s.processor.LikeStatus(
				ctx,
				ap.NewLikeActivity(object.Object.ActivityObject),
			)
			if err != nil {
				return err
			}
			return s.store.Favourites().DeleteByID(ctx, favourite.ID)
		}, nil
	default:
		return nil, errors.New("unsupported Activity type")
	}
}

func (s *federationService) processDelete(object *domain.ObjectOrLink) (worker.Job, error) {
	switch object.Object.Type {
	case "Note":
		return func(ctx context.Context) error {
			status, err := s.processor.LookupStatus(
				ctx,
				ap.NewNote(object.Object.ActivityObject),
			)
			if err != nil {
				return err
			}
			return s.store.Statuses().DeleteByID(ctx, status.ID)
		}, nil
	default:
		return nil, errors.New("unsupported Activity type")
	}
}

func (s *federationService) LookupRemoteActor(
	ctx context.Context,
	handle *util.HandleInfo,
) (*domain.Actor, error) {
	if handle == nil {
		return nil, errors.New("handle is required")
	}
	if handle.Domain == "" {
		return nil, errors.New("handle domain is required")
	}
	account := fmt.Sprintf("%s@%s", handle.Username, handle.Domain)
	resource := fmt.Sprintf("acct:%s", account)

	client := http.Client{Timeout: 5 * time.Second}
	webfingerURL := fmt.Sprintf("http://%s/.well-known/webfinger", handle.Domain)
	log.Printf("lookup_remote: requesting WebFinger for %s at %s", account, webfingerURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, webfingerURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("resource", resource)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("lookup_remote: WebFinger request failed status=%s", resp.Status)
		return nil, fmt.Errorf("webfinger lookup failed: %s", resp.Status)
	}
	var webfinger api.WebFingerResponse
	if err := json.NewDecoder(resp.Body).Decode(&webfinger); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	actorURL, err := selectActorLink(webfinger.Links)
	if err != nil {
		return nil, err
	}

	log.Printf("lookup_remote: fetching actor document from %s", actorURL)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, actorURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("lookup_remote: actor request failed status=%s", resp.Status)
		return nil, fmt.Errorf("actor lookup failed: %s", resp.Status)
	}

	var actor domain.Actor
	if err := json.NewDecoder(resp.Body).Decode(&actor); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &actor, nil
}

func selectActorLink(links []api.WebFingerLink) (string, error) {
	for _, link := range links {
		if link.Rel == "self" && link.Href != "" {
			if link.Type == "" || strings.Contains(link.Type, "activity+json") {
				return link.Href, nil
			}
		}
	}
	for _, link := range links {
		if link.Href != "" {
			return link.Href, nil
		}
	}
	return "", errors.New("webfinger response missing actor link")
}
