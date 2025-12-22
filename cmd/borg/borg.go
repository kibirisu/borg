package main

import (
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	conf := config.GetConfig()
	s := server.New(conf)
	panic(s.ListenAndServe())
}
