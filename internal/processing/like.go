package processing

import (
	"context"
	"errors"
	"log"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) LikeStatus(
	ctx context.Context,
	activity ap.LikeActivitier,
	targetAccountID xid.ID,
) (*xid.ID, error) {
	uri := activity.GetURI()
	log.Printf("[Processing] [Like] processing Activity with ID=%s", uri)
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	favourite, err := p.store.Favourites().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		log.Printf("[Processing] [Like] processing Actor with ID=%s", activityData.Actor.GetURI())
		accountID, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return nil, err
		}
		log.Printf("[Processing] [Like] processing Status with ID=%s", activityData.Object.GetURI())
		statusID, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return nil, err
		}
		favourite, err = p.store.Favourites().AddLike(ctx, db.AddLikeParams{
			ID:              xid.New(),
			Uri:             uri,
			AccountID:       *accountID,
			AccountUri:      activityData.Actor.GetURI(),
			StatusID:        *statusID,
			TargetAccountID: targetAccountID,
			StatusUri:       activityData.Object.GetURI(),
		})
		return &favourite.ID, err
	}
	return &favourite.ID, err
}
