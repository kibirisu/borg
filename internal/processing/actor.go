package processing

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/util"
)

func (p *processor) LookupActor(ctx context.Context, object ap.Actorer) (*xid.ID, error) {
	uri := object.GetURI()
	log.Printf("[Processing] [Actor] processing Object with ID=%s", uri)
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	account, err := p.store.Accounts().GetByURI(ctx, uri)
	if err != nil {
		object, err := p.client.Get(ctx, uri)
		if err != nil {
			return nil, err
		}
		actorData := ap.NewActor(object).GetObject()
		id := xid.New()
		account, err = p.store.Accounts().Create(ctx, db.CreateActorParams{
			ID:       id,
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
			FollowersUri: actorData.Followers,
			FollowingUri: actorData.Following,
		})
		return &account.ID, err
	}
	return &account.ID, err
}

// FetchAndStoreAccount implements Processor.
func (p *processor) FetchAndStoreAccount(
	ctx context.Context,
	username, domain string,
) (account db.Account, err error) {
	uri, err := p.getActorURI(ctx, username, domain)
	if err != nil {
		return
	}
	obj, err := p.client.Get(ctx, uri)
	if err != nil {
		return
	}
	actor := ap.NewActor(obj).GetObject()
	return p.store.Accounts().Add(ctx, db.AddAccountParams{
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
	})
}

func (p *processor) getActorURI(
	ctx context.Context,
	username, domain string,
) (uri string, err error) {
	webfinger, err := p.client.Webfinger(ctx, util.BuildWebfingerURL(username, domain))
	if err != nil {
		return
	}
	if len(webfinger.Links) != 1 {
		// we actually can process incoming webfinger, we decide to not
		return uri, errors.New("dealing with webfinger too advanced to understand")
	}
	return webfinger.Links[0].Href, nil
}
