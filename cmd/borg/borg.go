package main

import (
	"context"
	"fmt"

	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/server"
)

func main() {
	ctx := context.Background()
	conf := config.GetConfig()
	fmt.Println("--- Config Test ---")
    fmt.Printf("Env:  %s\n", conf.AppEnv)
    fmt.Printf("Host: %s\n", conf.ListenHost)
    fmt.Printf("Port: %s\n", conf.ListenPort)
    fmt.Printf("DB:   %s\n", conf.DatabaseURL)
    fmt.Println("-------------------")
	s := server.New(ctx, conf)
	panic(s.ListenAndServe())
}
