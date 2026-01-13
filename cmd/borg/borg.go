package main

import (
	"context"
	"log"

	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()
	log.Println("--- Config Test ---")
	log.Printf("Env:  %s\n", conf.AppEnv)
	log.Printf("Host: %s\n", conf.ListenHost)
	log.Printf("Port: %s\n", conf.ListenPort)
	log.Printf("DB:   %s\n", conf.DatabaseURL)
	log.Println("-------------------")
	s := server.New(ctx, conf)
	panic(s.ListenAndServe())
}
