package service

import (
	"github.com/kibirisu/borg/internal/config"
	proc "github.com/kibirisu/borg/internal/processing"
	"github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/transport"
)

type Container struct {
	App        AppService
	Federation FederationService
}

func NewContainer(conf *config.Config) *Container {
	store := repository.New(conf.DatabaseURL)
	proc := proc.New(store, transport.New())
	return &Container{
		App:        &appService{store, conf},
		Federation: &federationService{store, proc},
	}
}
