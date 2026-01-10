package processing

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) AnnounceStatus(ctx context.Context, activity ap.AnnounceActivitier) error {
	activityData := activity.GetObject()
	actor, err := p.LookupActor(ctx, activityData.Actor)
	if err != nil {
		return err
	}
	announcedStatus, err := p.LookupStatus(ctx, activityData.Object)
	if err != nil {
		return err
	}
	_, err = p.store.Statuses().Create(ctx, db.CreateStatusParams{
		Uri:       activityData.ID,
		AccountID: actor.ID,
		ReblogOfID: sql.NullInt32{
			Int32: announcedStatus.ID,
			Valid: true,
		},
	})
	return err
}
