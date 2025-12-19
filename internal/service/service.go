package service

import (
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/repository"
)

type Container struct {
	App        AppService
	Federation FederationService
}

func NewContainer(conf *config.Config) *Container {
	store := repository.NewStore(conf.DatabaseURL)
	return &Container{
		App:        NewAppService(store, conf),
		Federation: NewFederationService(store),
	}
}
