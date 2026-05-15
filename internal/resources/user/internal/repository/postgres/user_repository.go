package postgres

import (
	"context"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{
		q: q,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Name:         params.Name,
	})
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}
