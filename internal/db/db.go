package db

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"time"
)

var (
	//go:embed migrations/*
	embedMigrations embed.FS
)

func NewDb(dbStr string) (*pgxpool.Pool, error) {
	var dbPool *pgxpool.Pool
	var err error

	for attempts := 1; attempts <= 10; attempts++ {
		dbPool, err = pgxpool.New(context.Background(), dbStr)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: unable to create connection pool: %v\n", attempts, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool after 10 attempts: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	for attempts := 1; attempts <= 10; attempts++ {
		db := stdlib.OpenDBFromPool(dbPool)
		if err := goose.Up(db, "migrations"); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if err := db.Close(); err != nil {
			panic(err)
		}
		break
	}

	return dbPool, nil
}
