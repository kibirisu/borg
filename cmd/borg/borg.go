package main

import (
	"context"
	"log"

	"borg/internal/config"
	"borg/internal/domain"
	"borg/internal/server"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()

	ds, err := domain.NewDataStore(ctx, conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(conf.ListenPort, ds)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
