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

func (p *processor) AddActor(ctx context.Context, actor ap.Actorer) error {
	switch actor.GetValueType() {
	case ap.LinkType:
		uri := actor.GetLink()
		if util.ExtractDomainFromURI(uri) == p.conf.Address {
			return errors.New("attempted to fetch local account")
		}
		obj, err := p.client.Get(ctx, uri)
		if err != nil {
			return err
		}
		raw := actor.GetRaw()
		*raw = *obj
	case ap.ObjectType:
		uri := actor.GetRaw().Object.ID
		if util.ExtractDomainFromURI(uri) == p.conf.Address {
			return errors.New("attempted to fetch local account")
		}
	case ap.InvalidType:
		fallthrough
	case ap.NullType:
		fallthrough
	default:
		return errors.New("activity actor is missing")
	}
	obj := actor.GetObject()
	_, err := p.store.Accounts().Add(ctx, db.AddAccountParams{
		ID:       xid.New(),
		Username: obj.PreferredUsername,
		Uri:      obj.ID,
		Domain: sql.NullString{
			String: util.ExtractDomainFromURI(obj.ID),
			Valid:  true,
		},
		InboxUri:     obj.Inbox,
		OutboxUri:    obj.Outbox,
		FollowersUri: obj.Followers,
		FollowingUri: obj.Following,
		Url:          ":3",
	})
	return err
}

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
		Url:          ":3", // webfinger may provide url btw
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
		return uri, errors.New("dealing with webfinger too advanced to understand")
	}
	return webfinger.Links[0].Href, nil
}
