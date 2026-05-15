package database

import (
	"context"
	"database/sql"
	"log"
	"social-media-backend/internal/env"
	"social-media-backend/migrations"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDb() *sql.DB {
	db, err := sql.Open("pgx", env.Config.POSTGRES_CONNECTION)
	if err != nil {
		log.Fatalf("Falied to open db: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed To Ping Database: %v", err)
	}

	if err := migrations.MigrateFS(db, migrations.FS, "."); err != nil {
		log.Fatalf("Failed To Run DB Migrations: %v", err)
	}
	return db
}
