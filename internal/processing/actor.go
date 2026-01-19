package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/util"
	"github.com/rs/xid"
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
			Username:    actorData.PreferredUsername,
			Uri:         actorData.ID,
			DisplayName: sql.NullString{},
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
	object, err := p.fetchActor(ctx, username, domain)
	if err != nil {
		return db.Account{}, err
	}
	actor := ap.NewActor(object).GetObject()
	return p.store.Accounts().AddAccount(ctx, db.AddAccountParams{
		ID:       xid.New(),
		Username: actor.PreferredUsername,
		Uri:      actor.ID,
		Domain: sql.NullString{
			String: domain,
			Valid:  true,
		},
		InboxUri:     actor.Inbox,
		OutboxUri:    actor.Outbox,
		FollowersUri: actor.Followers,
		FollowingUri: actor.Following,
		Url:          ":3",
	})
}

func (p *processor) fetchActor(
	ctx context.Context,
	username, domain string,
) (*domain.ObjectOrLink, error) {
	webfinger, err := p.client.Webfinger(ctx, util.BuildWebfingerURL(username, domain))
	if err != nil {
		return nil, err
	}
	if len(webfinger.Links) != 1 {
		return nil, errors.New("dealing with webfinger too advanced to understand")
	}
	return p.client.Get(ctx, webfinger.Links[0].Href)
}
