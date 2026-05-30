package database

import (
	"context"
	"log"
	"social-media-backend/internal/env"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewDb() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, env.Config.POSTGRES_CONNECTION)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	stdDB := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := stdDB.Close(); err != nil {
			log.Printf("error while closing stdDB: %v\n", err)
		}
	}()

	// if err := migrations.MigrateFS(stdDB, migrations.FS, "."); err != nil {
	// 	log.Fatalf("Failed migrations: %v", err)
	// }

	return pool
}
