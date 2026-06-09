package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func NewDb(connectionStr string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(connectionStr)
	if err != nil {
		log.Fatalf("Failed to parse pool config: %v", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
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
