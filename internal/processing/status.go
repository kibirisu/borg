package processing

import (
	"context"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Status interface {
	Get(context.Context) (db.Status, error)
}

type status struct {
	object ap.Note
	store  repo.Store
	client transport.Client
}

var _ Status = (*status)(nil)

// Get implements Status.
func (s *status) Get(context.Context) (db.Status, error) {
	panic("unimplemented")
}
