package processing

import (
	"context"
	"errors"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) AnnounceStatus(
	ctx context.Context,
	activity ap.AnnounceActivitier,
) (db.Status, error) {
	uri := activity.GetURI()
	if uri == "" {
		return db.Status{}, errors.New("invalid object")
	}
	status, err := p.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		actorID, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return status, err
		}
		reblogOfID, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return status, err
		}
		return p.store.Statuses().Create(ctx, db.CreateStatusParams{
			AccountID:  *actorID,
			ReblogOfID: reblogOfID,
			ID:         xid.New(),
			AccountUri: activityData.Actor.GetURI(),
			Uri:        uri,
		})
	}
	return status, nil
}
