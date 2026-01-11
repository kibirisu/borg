package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)
func pop[T any](s *[]T) T {
	lastIdx := len(*s) - 1
	element := (*s)[lastIdx]
	*s = (*s)[:lastIdx]
	return element
}

func (p *processor) FetchActorCollectionPage(ctx context.Context, collectionUri string) (ap.ActorCollectionPager, error) {
	col, err := p.client.Get(ctx, collectionUri)
	if err != nil {
		return ap.NewActorCollectionPage(nil), err
	}
	fetchedCollection := ap.NewActorCollection(col)
	collectionData := fetchedCollection.GetObject()
	pageUri := collectionData.First.GetURI()
	// we assume only one page is used
	colP, err := p.client.Get(ctx, pageUri)
	if err != nil {
		return ap.NewActorCollectionPage(nil), err
	}
	fetchedCollectionPage := ap.NewActorCollectionPage(colP)
	return fetchedCollectionPage, nil
}
func (p *processor) LookupFollowers(ctx context.Context, account db.Account, actor ap.Actor) error {
	followersUri := actor.Followers
	page, err := p.FetchActorCollectionPage(ctx, followersUri)
	if err != nil {
		return err
	}
	for _, follower := range page.GetObject().Items {
		// this will recursively call lookupfollowers for other accounts
		dbFollower, err := p.LookupActor(ctx, follower)
		if err != nil {
			return err
		}
		params := db.CreateFollowParams {
			AccountID: dbFollower.ID,
			TargetAccountID: account.ID,
		}
		_, err = p.store.Follows().Create(ctx, params)
		if err != nil {
			return err
		}
	}

	return nil
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
			return account, err
		}
		err = p.LookupFollowers(ctx, account, actorData)
		if err != nil {
			return account, err
		}
	}
	return account, nil
}
