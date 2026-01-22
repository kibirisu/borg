package processing

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func (p *processor) LookupStatus(ctx context.Context, object ap.Noter) (*xid.ID, error) {
	uri := object.GetURI()
	log.Printf("[Processing] [Status] processing Object with ID=%s", uri)
	if uri == "" {
		return nil, errors.New("invalid object")
	}
	status, err := p.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		object, err := p.client.Get(ctx, uri)
		if err != nil {
			return nil, err
		}
		statusData := ap.NewNote(object).GetObject()
		log.Printf(
			"[Processing] [Status] processing Actor with ID=%s",
			statusData.AttributedTo.GetURI(),
		)
		accountID, err := p.LookupActor(ctx, statusData.AttributedTo)
		if err != nil {
			return nil, err
		}
		var inReplyToID *xid.ID
		if statusData.InReplyTo.GetRaw() != nil {
			inReplyToID, err = p.LookupStatus(ctx, statusData.InReplyTo)
			if err != nil {
				return nil, err
			}
		}
		status, err = p.store.Statuses().Add(ctx, db.AddStatusParams{
			ID:        xid.New(),
			AccountID: *accountID,
			Content: sql.NullString{
				String: statusData.Content,
				Valid:  true,
			},
			InReplyToID: inReplyToID,
			Uri:         uri,
		})
		if err != nil {
			return nil, err
		}
	}
	return &status.ID, err
}
