package processing

import (
	"github.com/kibirisu/borg/internal/ap"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Processor interface {
	Actor(ap.Actorer) Actor
	Status(ap.Noter) Status
}

type processor struct {
	store  repo.Store
	client transport.Client
}

var _ Processor = (*processor)(nil)

func NewProcessor(store repo.Store, client transport.Client) Processor {
	return &processor{store, client}
}

// Actor implements Processor.
func (p *processor) Actor(object ap.Actorer) Actor {
	return &actor{object, p.store, p.client}
}

// Status implements Processor.
func (p *processor) Status(object ap.Noter) Status {
	return &status{object, p.store, p.client}
}
