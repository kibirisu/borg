package processing

import (
	"context"
	"errors"

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

func (p *processor) FollowStatus(
	ctx context.Context,
	activity ap.FollowActivitier,
) (db.Follow, error) {
	uri := activity.GetURI()
	if uri == "" {
		return db.Follow{}, errors.New("invalid object")
	}
	follow, err := p.store.Follows().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		followerAccount, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return follow, err
		}
		followedActor, err := p.LookupActor(ctx, activityData.Object)
		if err != nil {
			return follow, err
		}
		DBfollow, err := p.store.Follows().Create(ctx, db.CreateFollowParams{
			Uri:             "/users/" + followedActor.Username + "/followers/" + followerAccount.Username, // TODO
			AccountID:       followerAccount.ID,
			TargetAccountID: followedActor.ID,
		})
		return *DBfollow, err
	}
	return follow, nil
}
