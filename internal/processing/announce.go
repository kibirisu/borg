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
) (*xid.ID, error) {
	uri := activity.GetURI()
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	status, err := p.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		actorID, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return nil, err
		}
		reblogOfID, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return nil, err
		}
		status, err = p.store.Statuses().Create(ctx, db.CreateStatusParams{
			AccountID:  *actorID,
			ReblogOfID: reblogOfID,
			ID:         xid.New(),
			AccountUri: activityData.Actor.GetURI(),
			Uri:        uri,
		})
		if err != nil {
			return nil, err
		}
	}
	return &status.ID, err
}
