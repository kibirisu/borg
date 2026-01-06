package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

type Actor interface {
	Get(context.Context) (db.Account, error)
}

type actor struct {
	object    ap.Actor
	processor *processor
}

var _ Actor = (*actor)(nil)

// Get implements Actor.
func (a *actor) Get(ctx context.Context) (db.Account, error) {
	account, err := a.processor.store.Accounts().GetByURI(ctx, a.object.ID)
	_, _ = account, err
	panic("")
}
