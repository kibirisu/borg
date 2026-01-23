package main

import (
	"context"

	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()
	s := server.New(ctx, conf)
	panic(s.ListenAndServe())
}
