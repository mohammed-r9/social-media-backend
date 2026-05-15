package postgres

import (
	"context"
	"social-media-backend/internal/adapters/sqlc/repo"
	"social-media-backend/internal/domain/user"
)

type UserRepository struct {
	q *repo.Queries
}

// function is unfinished
func (u *UserRepository) CreateUser(ctx context.Context, params user.CreateUserParams) {
	user, err := u.q.CreateUser(ctx, repo.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: []byte(params.PassowrdHash), // maybe temp
		Name:         params.Name,
	})
	_, _ = user, err
}
