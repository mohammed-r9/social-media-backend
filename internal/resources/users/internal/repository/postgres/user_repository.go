package postgres

import (
	"context"
	"errors"
	"fmt"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
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
	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Name:         params.Name,
	})

	if err != nil {
		if pqErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pqErr.Code == "23505" {
				return domain.User{}, domain.ErrEmailAlreadyTaken
			}
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return domain.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
