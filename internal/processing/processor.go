package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Processor interface {
	LookupActor(context.Context, ap.Actorer) (db.Account, error)
	LookupStatus(context.Context, ap.Noter) (db.Status, error)
	AnnounceStatus(context.Context, ap.AnnounceActivitier) (db.Status, error)
	AcceptFollow(context.Context, ap.FollowActivitier) error
	AcceptLike(context.Context, ap.LikeActivitier) (db.Favourite, error)
}

type processor struct {
	store  repo.Store
	client transport.Client
}

var _ Processor = (*processor)(nil)

func New(store repo.Store, client transport.Client) Processor {
	return &processor{store, client}
}
