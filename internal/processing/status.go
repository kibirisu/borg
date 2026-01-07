package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Status interface {
	Get(context.Context) (db.Status, error)
}

type status struct {
	object ap.Noter
	store  repo.Store
	client transport.Client
}

var _ Status = (*status)(nil)

func newStatus(object ap.Noter, store repo.Store, client transport.Client) Status {
	return &status{object, store, client}
}

// Get implements Status.
func (s *status) Get(ctx context.Context) (db.Status, error) {
	uri := s.object.GetURI()
	if uri == "" {
		return db.Status{}, errors.New("invalid object")
	}
	status, err := s.store.Statuses().GetByURI(ctx, uri)
	if err != nil {
		object, err := s.client.Get(uri)
		if err != nil {
			return status, err
		}
		fetchedStatus := ap.NewNote(object)
		statusData := fetchedStatus.GetObject()
		account, err := newActor(statusData.AttributedTo, s.store, s.client).Get(ctx)
		if err != nil {
			return status, err
		}
		_ = account
		inReplyToID := sql.NullInt32{}
		if statusData.InReplyTo.GetRaw() != nil {
			parentStatus, err := newStatus(statusData.InReplyTo, s.store, s.client).Get(ctx)
			if err != nil {
				return status, err
			}
			inReplyToID = sql.NullInt32{
				Int32: parentStatus.ID,
				Valid: true,
			}
		}
		status, err = s.store.Statuses().Create(ctx, db.CreateStatusParams{
			Uri:         uri,
			Url:         "nope",
			Local:       sql.NullBool{},
			Content:     statusData.Content,
			AccountID:   account.ID,
			InReplyToID: inReplyToID,
			ReblogOfID:  sql.NullInt32{},
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}
