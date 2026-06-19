package repo

import (
	"social-media-backend/internal/adapters/sqlc/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool, q *db.Queries) *PostgresRepo {
	return &PostgresRepo{
		queries: q,
		db:      db,
	}
}
