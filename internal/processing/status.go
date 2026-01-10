package processing

import (
	"context"
	"database/sql"
	"errors"

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
		inReplyToID := sql.NullInt32{}
		if statusData.InReplyTo.GetRaw() != nil {
			parentStatus, err := p.LookupStatus(ctx, statusData.InReplyTo)
			if err != nil {
				return status, err
			}
			inReplyToID = sql.NullInt32{
				Int32: parentStatus.ID,
				Valid: true,
			}
		}
		status, err = p.store.Statuses().Create(ctx, db.CreateStatusParams{
			Uri:         uri,
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
