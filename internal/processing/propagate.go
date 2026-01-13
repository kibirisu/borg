package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
)

func (p *processor) Propagate(ctx context.Context, activity ap.Activiter[any]) error {
	senderActor := ap.NewActor(activity.GetRaw().Object.ActivityActor)
	_ = senderActor
	// p.store.Accounts().GetFollowers(ctx, 1)
	panic("")
}
