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
func (p *processor) LookupFollowers(ctx context.Context, account ap.Actorer) ([]db.Account, []db.Follow, error) {
	type ActorOrigin struct {
		Actor		ap.Actorer
		FollowerOf	int32
	}
	accounts := make([]db.Account, 0)
	follows := make([]db.Follow, 0)
	toProcess := make([]ActorOrigin, 0)
	toProcess = append(toProcess, ActorOrigin{account, -1})
	for len(toProcess) > 0 {
		cur := pop(&toProcess)
		actorURI := cur.Actor.GetURI()
		// skip if already in db (so theoretically already resolved)
		_, err := p.store.Accounts().GetByURI(ctx, actorURI)
		if err == nil {
			continue
		}
		// fetch actor that hasnt yet been processed
		acc, err := p.LookupActor(ctx, cur.Actor)
		if err != nil {
			continue
		}
		accounts = append(accounts, acc)
		// if it wasnt root account (first one searched) then add the follow to parent
		if cur.FollowerOf != -1 {
			param := db.CreateFollowParams {
				AccountID:			acc.ID,
				TargetAccountID:	cur.FollowerOf, 
			}
			follow, err := p.store.Follows().Create(ctx, param)
			if err == nil {
				follows = append(follows, *follow)
			}
		}
		// add this account followers as well
		followersUri := cur.Actor.GetObject().Followers
		page, err := p.FetchActorCollectionPage(ctx, followersUri)
		if err != nil {
			continue
		}
		for _, follower := range page.GetObject().Items {
			toProcess = append(toProcess, ActorOrigin{
				Actor:      follower, 
				FollowerOf: acc.ID,
			})
		}
	}

	return accounts, follows, nil
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
	}
	return account, nil
}
