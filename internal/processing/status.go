package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) LookupStatus(ctx context.Context, object ap.Noter) (db.Status, error) {
	uri := object.GetURI()
	if uri == "" {
		return db.Status{}, errors.New("invalid object")
	}
	status, err := p.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		object, err := p.client.Get(ctx, uri)
		if err != nil {
			return status, err
		}
		fetchedStatus := ap.NewNote(object)
		statusData := fetchedStatus.GetObject()
		account, err := p.LookupActor(ctx, statusData.AttributedTo)
		if err != nil {
			return status, err
		}
		var inReplyToID *xid.ID
		if statusData.InReplyTo.GetRaw() != nil {
			parentStatus, err := p.LookupStatus(ctx, statusData.InReplyTo)
			if err != nil {
				return status, err
			}
			inReplyToID = &parentStatus.ID
		}
		status, err = p.store.Statuses().Create(ctx, db.CreateStatusParams{
			Url:         "nope",
			Local:       sql.NullBool{},
			Content:     statusData.Content,
			AccountID:   account.ID,
			InReplyToID: inReplyToID,
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}

func (p *processor) DistributeStatus(
	ctx context.Context,
	activity ap.CreateActivitier,
	actorID xid.ID,
) error {
	inboxes, err := p.store.Accounts().GetAccountRemoteFollowerInboxes(ctx, actorID)
	if err != nil {
		return err
	}
	object := activity.GetRaw().Object
	for _, inbox := range inboxes {
		if err = p.client.Post(ctx, inbox, object); err != nil {
			return err
		}
	}
	return nil
}
