package processing

import (
	"context"
	"log"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/domain"
)

func (p *processor) DistributeObject(
	ctx context.Context,
	object *domain.Object,
	actorID xid.ID,
) error {
	inboxes, err := p.store.Accounts().GetAccountRemoteFollowerInboxes(ctx, actorID)
	if err != nil {
		return err
	}
	for _, inbox := range inboxes {
		if err = p.client.Post(ctx, inbox, object); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (p *processor) SendObject(
	ctx context.Context,
	object *domain.Object,
	receiverID xid.ID,
) error {
	inbox, err := p.store.Accounts().GetAccountInbox(ctx, receiverID)
	if err != nil {
		return err
	}
	return p.client.Post(ctx, inbox, object)
}

