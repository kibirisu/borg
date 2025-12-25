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

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/util"
)

type FederationService interface {
	GetLocalActor(context.Context, string) (*domain.Actor, error)
	CreateActor(context.Context, db.CreateActorParams) (*db.Account, error)
	LookupRemoteActor(context.Context, *util.HandleInfo) (*domain.Actor, error)
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

	log.Printf("lookup_remote: actor %s retrieved successfully", actor.ID)
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
