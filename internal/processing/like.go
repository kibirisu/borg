package processing

import (
	"context"
	"errors"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) LikeStatus(
	ctx context.Context,
	activity ap.LikeActivitier,
) (*xid.ID, error) {
	uri := activity.GetURI()
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	favourite, err := p.store.Favourites().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		accountID, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return nil, err
		}
		statusID, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return nil, err
		}
		favourite, err = p.store.Favourites().Create(ctx, db.CreateFavouriteParams{
			AccountID: *accountID,
			StatusID:  *statusID,
		})
		if err != nil {
			return nil, err
		}
	}
	return &favourite.ID, nil
}
