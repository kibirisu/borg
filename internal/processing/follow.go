package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) AcceptFollow(ctx context.Context, activity ap.FollowActivitier) error {
	activityData := activity.GetObject()
	localAccount, err := p.store.Accounts().GetByURI(ctx, activityData.Object.GetURI())
	if err != nil {
		return err
	}
	remoteAccount, err := p.LookupActor(ctx, activityData.Actor)
	if err != nil {
		return err
	}
	_, err = p.store.Follows().Create(ctx, db.CreateFollowParams{
		Uri:             activityData.ID,
		AccountID:       remoteAccount.ID,
		TargetAccountID: localAccount.ID,
	})
	if err != nil {
		return err
	}

	accept := ap.NewAcceptActivity(nil)
	accept.SetObject(ap.Activity[ap.Activity[ap.Actor]]{
		ID:     "TODO",
		Type:   "Accept",
		Actor:  activityData.Object,
		Object: activity,
	})
	return p.client.Post(ctx, remoteAccount.InboxUri, accept.GetRaw().Object)
}
