package processing

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/ap"
)

func (p *processor) Propagate(ctx context.Context, activity ap.Activiter[any]) error {
	senderActor := ap.NewActor(activity.GetRaw().Object.ActivityActor)
	_ = senderActor
	p.store.Accounts().GetFollowers(ctx, xid.New())
	panic("")
}
