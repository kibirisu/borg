package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) AcceptLike(ctx context.Context, activity ap.LikeActivitier) error {
	activityData := activity.GetObject()
	likedAccount, err := p.LookupActor(ctx, activityData.Actor)
	likedPost, err := p.LookupStatus(ctx, activityData.Object)
	if err != nil {
		return err
	}
	_, err = p.store.Favourites().Create(ctx, db.CreateFavouriteParams{
		AccountID:  likedAccount.ID,
		StatusID:	likedPost.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
