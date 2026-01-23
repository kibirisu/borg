package processing

import (
	"context"
	"errors"
	"log"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) AnnounceStatus(
	ctx context.Context,
	activity ap.AnnounceActivitier,
) (*xid.ID, error) {
	uri := activity.GetURI()
	log.Printf("[Processing] [Announce] processing Activity with ID=%s", uri)
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	status, err := p.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		activityData := activity.GetObject()
		log.Printf(
			"[Processing] [Announce] processing Actor with ID=%s",
			activityData.Actor.GetURI(),
		)
		actorID, err := p.LookupActor(ctx, activityData.Actor)
		if err != nil {
			return nil, err
		}
		log.Printf(
			"[Processing] [Announce] processing Status with ID=%s",
			activityData.Object.GetURI(),
		)
		reblogOfID, err := p.LookupStatus(ctx, activityData.Object)
		if err != nil {
			return nil, err
		}
		status, err = p.store.Statuses().AddReblog(ctx, db.AddReblogParams{
			ID:         xid.New(),
			ReblogUri:  uri,
			AccountID:  *actorID,
			ReblogOfID: reblogOfID,
		})
		return &status.ID, err
	}
	return &status.ID, err
}
