package postgres

import "social-media-backend/internal/adapters/sqlc/db"

type PostgresRepo struct {
	q *db.Queries
}

func NewPostgresRepository(q *db.Queries) *PostgresRepo {
	return &PostgresRepo{
		q: q,
	}
}
