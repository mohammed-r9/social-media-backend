package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

var _ UserRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
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
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Phone:        utils.NullStringToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	if email == "" {
		return domain.User{}, domain.ErrInvalidUserEmail
	}

	user, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Phone:        utils.NullStringToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}
