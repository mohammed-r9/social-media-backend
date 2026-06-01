package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"
	"social-media-backend/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	UpdateUserPassword(ctx context.Context, params domain.UpdatePasswordParams) error
	VerifyUserEmail(ctx context.Context, userID uuid.UUID) error
}

var _ UserRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Username:     params.Username,
	})
	if err != nil {
		if pqErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pqErr.Code == "23505" {
				return domain.User{}, domain.ErrEmailAlreadyTaken
			}
		}

		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	_, err = qtx.CreateProfile(ctx, db.CreateProfileParams{
		UserID:      user.ID,
		DisplayName: params.Name,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("create profile: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("commit tx: %w", err)
	}

	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
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
		Username:     user.Username,
		Phone:        utils.NullStringToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	if userID == uuid.Nil {
		return domain.User{}, domain.ErrInvalidUserID
	}

	user, err := r.q.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	return domain.User{
		ID:    user.ID,
		Email: user.Email,
		Profile: domain.Profile{
			DisplayName:    utils.NullStringToString(user.DisplayName),
			Bio:            utils.NullStringToString(user.Bio),
			AvatarUrl:      utils.NullStringToString(user.AvatarKey),
			Website:        utils.NullStringToString(user.Website),
			FollowerCount:  utils.PgBigIntToInt64(user.FollowersCount),
			FollowingCount: utils.PgBigIntToInt64(user.FollowingCount),
			PostCount:      utils.PgBigIntToInt64(user.PostsCount),
		},
		Username:     user.Username,
		Phone:        utils.NullStringToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) UpdateUserPassword(ctx context.Context, params domain.UpdatePasswordParams) error {
	rowsAffected, err := r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           params.ID,
		PasswordHash: params.PasswordHash,
	})

	if rowsAffected == 0 {
		return domain.ErrNoRowsAffected
	}
	return err
}

func (r *PostgresRepo) VerifyUserEmail(ctx context.Context, userID uuid.UUID) error {
	rowsAffected, err := r.q.VerifyUserEmail(ctx, userID)
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return err
}
