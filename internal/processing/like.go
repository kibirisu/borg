package processing

import (
	"context"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) LikeStatus(
	ctx context.Context,
	activity ap.LikeActivitier,
) (db.Favourite, error) {
	uri := activity.GetURI()
	if uri == "" {
		return db.Favourite{}, errors.New("invalid object")
	}
	favourite, err := p.store.Favourites().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		likerAccount, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return favourite, err
		}
		likedPost, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return favourite, err
		}
		return p.store.Favourites().Create(ctx, db.CreateFavouriteParams{
			AccountID: likerAccount.ID,
			StatusID:  likedPost.ID,
		})
	}
	return favourite, nil
}
