package processing

import (
	"context"
	"database/sql"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/util"
)

func (p *processor) AcceptFollow(
	ctx context.Context,
	activity ap.FollowActivitier,
	targetAccountID xid.ID,
) error {
	activityData := activity.GetObject()
	followID := xid.New()
	inbox, err := p.store.Follows().AddByActorURI(ctx, db.AddFollowByActorURIParams{
		ID:              followID,
		Uri:             activityData.ID,
		AccountUri:      activityData.Actor.GetURI(),
		TargetAccountID: targetAccountID,
	})
	if err != nil {
		obj, err := p.client.Get(ctx, activityData.Actor.GetURI())
		if err != nil {
			return err
		}
		actor := ap.NewActor(obj).GetObject()

		res, err := p.store.WithTX(ctx, func(ctx context.Context, s repo.Store) (any, error) {
			accountID := xid.New()
			_, err := s.Accounts().Add(ctx, db.AddAccountParams{
				ID:       accountID,
				Username: actor.PreferredUsername,
				Uri:      actor.ID,
				Domain: sql.NullString{
					String: util.ExtractDomainFromURI(actor.ID),
					Valid:  true,
				},
				InboxUri:     actor.Inbox,
				OutboxUri:    actor.Outbox,
				FollowersUri: actor.Followers,
				FollowingUri: actor.Following,
				Url:          ":3",
			})
			if err != nil {
				return nil, err
			}
			err = s.Follows().Create(ctx, db.CreateFollowParams{
				ID:              followID,
				Uri:             activityData.ID,
				AccountID:       accountID,
				TargetAccountID: targetAccountID,
			})
			return actor.Inbox, err
		})
		if err != nil {
			return err
		}
		inbox = res.(string)
	}

	accept := ap.NewEmptyAcceptActivity().WithObject(ap.Activity[ap.Activity[ap.Actor]]{
		ID:     "nope",
		Type:   "Accept",
		Actor:  activityData.Object,
		Object: activity,
	})
	return p.client.Post(ctx, inbox, accept.GetRaw().Object)
}
