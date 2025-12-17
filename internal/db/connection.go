package db

import (
	"context"
	"database/sql"
	"embed"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func GetDB(ctx context.Context, url string) *Queries {
	pool, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
	}

	// Run database migrations
	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err = goose.Up(pool, "migrations"); err != nil {
		log.Fatal(err)
	}

	return New(pool)
}
