package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/util"
)

func (p *processor) LookupActor(ctx context.Context, object ap.Actorer) (db.Account, error) {
	uri := object.GetURI()
	if uri == "" {
		return db.Account{}, errors.New("invalid object")
	}
	account, err := p.store.Accounts().GetByURI(ctx, uri)
	if err != nil {
		object, err := p.client.Get(ctx, uri)
		if err != nil {
			return account, err
		}
		fetchedActor := ap.NewActor(object)
		actorData := fetchedActor.GetObject()
		account, err = p.store.Accounts().Create(ctx, db.CreateActorParams{
			Username: actorData.PreferredUsername,
			Uri:      actorData.ID,
			DisplayName: sql.NullString{
				String: actorData.Name,
				Valid:  true,
			},
			Domain: sql.NullString{
				String: util.ExtractDomainFromURI(uri),
				Valid:  true,
			},
			InboxUri:     actorData.Inbox,
			OutboxUri:    actorData.Outbox,
			Url:          "nope",
			FollowersUri: actorData.Followers,
			FollowingUri: actorData.Following,
		})
		if err != nil {
			return account, err
		}
	}
	return account, nil
}

// FetchAndStoreAccount implements Processor.
func (p *processor) FetchAndStoreAccount(
	ctx context.Context,
	username, domain string,
) (db.Account, error) {
	actor, err := p.fetchActor(ctx, username, domain)
	if err != nil {
		return db.Account{}, err
	}
	object := actor.GetObject()
	return p.store.Accounts().AddAccount(ctx, db.AddAccountParams{
		ID:       xid.New(),
		Username: object.PreferredUsername,
		Uri:      object.ID,
		Domain: sql.NullString{
			String: domain,
			Valid:  true,
		},
		InboxUri:     object.Inbox,
		OutboxUri:    object.Outbox,
		FollowersUri: object.Followers,
		FollowingUri: object.Following,
		Url:          ":3", // webfinger may provide url btw
	})
}

func (p *processor) fetchActor(
	ctx context.Context,
	username, domain string,
) (ap.Actorer, error) {
	webfinger, err := p.client.Webfinger(ctx, util.BuildWebfingerURL(username, domain))
	if err != nil {
		return nil, err
	}
	if len(webfinger.Links) != 1 {
		return nil, errors.New("dealing with webfinger too advanced to understand")
	}
	obj, err := p.client.Get(ctx, webfinger.Links[0].Href)
	return ap.NewActor(obj), err
}
