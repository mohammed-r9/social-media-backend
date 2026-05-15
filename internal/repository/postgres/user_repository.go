package postgres

import (
	"context"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain/user"
)

type UserRepository struct {
	q *db.Queries
}

// function is unfinished
func (u *UserRepository) CreateUser(ctx context.Context, params user.CreateUserParams) {
	user, err := u.q.CreateUser(ctx, db.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: params.PassowrdHash,
		Name:         params.Name,
	})
	_, _ = user, err
}
