package main

import (
	"context"
	"log"

	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/domain"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()

	ds, err := domain.NewDataStore(ctx, conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(conf, ds)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
