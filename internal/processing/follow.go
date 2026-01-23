package processing

import (
	"context"
	"database/sql"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/util"
)

func (p *processor) AcceptFollow(
	ctx context.Context,
	follow ap.FollowActivitier,
	targetAccountID xid.ID,
) error {
	activity := follow.GetObject()
	followID := xid.New()
	inbox, err := p.store.Follows().AddByActorURI(ctx, db.AddFollowByActorURIParams{
		ID:              followID,
		Uri:             activity.ID,
		AccountUri:      activity.Actor.GetURI(),
		TargetAccountID: targetAccountID,
	})
	if err != nil {
		if activity.Actor.GetValueType() != ap.ObjectType {
			obj, err := p.client.Get(ctx, activity.Actor.GetURI())
			if err != nil {
				return err
			}
			*activity.Actor.GetRaw() = *obj
		}
		actor := activity.Actor.GetObject()
		if err = p.store.Follows().AddWithActor(ctx, db.AddFollowWithActorParams{
			AccountID:  xid.New(),
			Username:   actor.PreferredUsername,
			AccountUri: actor.ID,
			DisplayName: sql.NullString{
				String: actor.Name,
				Valid:  true,
			},
			Domain: sql.NullString{
				String: util.ExtractDomainFromURI(activity.ID),
				Valid:  true,
			},
			Inbox:           actor.Inbox,
			Outbox:          actor.Outbox,
			Followers:       actor.Followers,
			Following:       actor.Following,
			FollowID:        followID,
			FollowUri:       activity.ID,
			TargetAccountID: targetAccountID,
		}); err != nil {
			return err
		}
		inbox = actor.Inbox
	}

	accept := ap.NewEmptyAcceptActivity().WithObject(ap.Activity[ap.Activity[ap.Actor]]{
		Type:   "Accept",
		Actor:  activity.Object,
		Object: follow,
	})
	return p.client.Post(ctx, inbox, accept.GetRaw().Object)
}
