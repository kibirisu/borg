package processing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/util"
)

func (p *processor) AddStatus(ctx context.Context, status ap.Noter) (*xid.ID, error) {
	switch status.GetValueType() {
	case ap.LinkType:
		uri := status.GetLink()
		if util.ExtractDomainFromURI(uri) == p.conf.Address {
			return nil, errors.New("attempted to fetch local status")
		}
		obj, err := p.client.Get(ctx, uri)
		if err != nil {
			return nil, err
		}
		raw := status.GetRaw()
		*raw = *obj
	case ap.ObjectType:
		uri := status.GetRaw().Object.ID
		if util.ExtractDomainFromURI(uri) == p.conf.Address {
			return nil, errors.New("attempted to fetch local status")
		}
	case ap.NullType:
		fallthrough
	case ap.InvalidType:
		fallthrough
	default:
		return nil, errors.New("activity status object is missing")
	}
	obj := status.GetObject()
	inReplyTo := obj.InReplyTo.GetRaw()
	if inReplyTo == nil {
		id := xid.New()
		err := p.store.Statuses().AddByActorURI(ctx, db.AddStatusByActorURIParams{
			ID:  id,
			Uri: obj.ID,
			Url: ":3",
			Content: sql.NullString{
				String: obj.Content,
				Valid:  true,
			},
			AccountUri: obj.AttributedTo.GetURI(),
		})
		if err != nil {
			if err = p.AddActor(ctx, obj.AttributedTo); err != nil {
				return nil, err
			}
			return nil, errors.New("unimplemented")
		}
		return &id, nil
	}
	return nil, errors.New("unimplemented")
}

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
			Url:   "nope",
			Local: sql.NullBool{},
			Content: sql.NullString{
				String: statusData.Content,
				Valid:  true,
			},
			AccountID:   account.ID,
			InReplyToID: inReplyToID,
		})
		if err != nil {
			return status, err
		}
	}
	return status, nil
}
