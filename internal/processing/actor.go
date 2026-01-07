package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Actor interface {
	Get(context.Context) (db.Account, error)
}

type actor struct {
	object ap.Actorer
	store  repo.Store
	client transport.Client
}

var _ Actor = (*actor)(nil)

func newActor(object ap.Actorer, store repo.Store, client transport.Client) Actor {
	return &actor{object, store, client}
}

// Get implements Actor.
func (a *actor) Get(ctx context.Context) (db.Account, error) {
	uri := a.object.GetURI()
	if uri == "" {
		return db.Account{}, errors.New("invalid object")
	}
	account, err := a.store.Accounts().GetByURI(ctx, uri)
	if err != nil {
		object, err := a.client.Get(uri)
		if err != nil {
			return account, err
		}
		fetchedActor := ap.NewActor(object)
		actorData := fetchedActor.GetObject()
		account, err := a.store.Accounts().Create(ctx, db.CreateActorParams{
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
