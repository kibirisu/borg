package processing

import (
	"github.com/kibirisu/borg/internal/ap"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Processor interface {
	Actor(ap.Actor) Actor
}

type processor struct {
	store  repo.Store
	client transport.Client
}

var _ Processor = (*processor)(nil)

// Actor implements Processor.
func (p *processor) Actor(object ap.Actor) Actor {
	return &actor{object, p}
}
