package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func GetDB(ctx context.Context, url string) (*Queries, error) {
	pool, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	// Run database migrations
	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err = goose.Up(pool, "migrations"); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to DB")
	return New(pool), nil
}
