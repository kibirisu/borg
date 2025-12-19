package main

import (
	// "context"

	"github.com/kibirisu/borg/internal/config"
	// "github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	// ctx := context.Background()
	conf := config.GetConfig()
	// _ = domain.NewDataStore(ctx, conf.DatabaseURL)
	s := server.New(conf)
	panic(s.ListenAndServe())
}
