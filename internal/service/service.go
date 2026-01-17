package service

import (
	"context"

	"github.com/kibirisu/borg/internal/config"
	proc "github.com/kibirisu/borg/internal/processing"
	"github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
	"github.com/kibirisu/borg/internal/util"
)

type Container struct {
	App        AppService
	Federation FederationService
}

func NewContainer(ctx context.Context, conf *config.Config) *Container {
	store := repository.New(ctx, conf.DatabaseURL)
	proc := proc.New(store, transport.New())
	builder := util.NewURIBuilder(conf.Address)
	return &Container{
		App:        &appService{store, proc, conf, builder},
		Federation: &federationService{store, proc},
	}
}
