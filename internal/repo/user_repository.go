package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/apperrors"
	"social-media-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	UpdateUserPassword(ctx context.Context, params domain.UpdatePasswordParams) (uuid.UUID, error)
	VerifyUserEmail(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	UpdateSelfProfile(ctx context.Context, params domain.UpdateProfileParams) (domain.Profile, error)
	UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarKey string) error
}

var _ UserRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		ID:           params.ID,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
		Username:     params.Username,
	})
	if err != nil {
		if pqErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pqErr.Code == "23505" {
				return domain.User{}, apperrors.EmailAlreadyTaken
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
		Phone:        textToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	if email == "" {
		return domain.User{}, apperrors.InvalidUserEmail
	}

	user, err := r.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, apperrors.UserNotFound
	}

	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		Phone:        textToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	if userID == uuid.Nil {
		return domain.User{}, apperrors.InvalidUserID
	}

	user, err := r.queries.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, apperrors.UserNotFound
	}

	return domain.User{
		ID:    user.ID,
		Email: user.Email,
		Profile: domain.Profile{
			DisplayName:    textToString(user.DisplayName),
			Bio:            textToString(user.Bio),
			AvatarKey:      textToString(user.AvatarKey),
			Website:        textToString(user.Website),
			FollowerCount:  int8ToInt64(user.FollowersCount),
			FollowingCount: int8ToInt64(user.FollowingCount),
			PostCount:      int8ToInt64(user.PostsCount),
		},
		Username:     user.Username,
		Phone:        textToString(user.Phone),
		PasswordHash: user.PasswordHash,
		IsSuspended:  user.IsSuspended,
		VerifiedAt:   user.VerifiedAt,
	}, nil
}

func (r *PostgresRepo) UpdateUserPassword(ctx context.Context, params domain.UpdatePasswordParams) (uuid.UUID, error) {
	userID, err := r.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           params.ID,
		PasswordHash: params.PasswordHash,
	})

	if userID == uuid.Nil {
		return uuid.Nil, apperrors.NoRowsAffected
	}

	return userID, err
}

func (r *PostgresRepo) VerifyUserEmail(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	userID, err := r.queries.VerifyUserEmail(ctx, userID)

	if userID == uuid.Nil {
		return uuid.Nil, apperrors.NoRowsAffected
	}

	return userID, err
}

func (r *PostgresRepo) UpdateSelfProfile(ctx context.Context, params domain.UpdateProfileParams) (domain.Profile, error) {
	if params.UserID == uuid.Nil {
		return domain.Profile{}, apperrors.InvalidUserID
	}

	profile, err := r.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		UserID:      params.UserID,
		Bio:         stringPtrToTex(params.Bio),
		DisplayName: stringPtrToTex(params.DisplayName),
		Website:     stringPtrToTex(params.Website),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {

			if errors.Is(err, pgx.ErrNoRows) {
				return domain.Profile{}, apperrors.ProfileNotFound
			}

			return domain.Profile{}, fmt.Errorf("update user profile: %w", err)
		}
	}

	return domain.Profile{
		DisplayName: profile.DisplayName,
		Bio:         textToString(profile.Bio),
		Website:     textToString(profile.Website),
		AvatarKey:   textToString(profile.AvatarKey),
	}, nil
}

func (r *PostgresRepo) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarKey string) error {
	if userID == uuid.Nil {
		return apperrors.InvalidUserID
	}
	_, err := r.queries.UpdateUserAvatar(ctx, db.UpdateUserAvatarParams{
		AvatarKey: stringToTex(avatarKey),
		UserID:    userID,
	})
	// TODO: add better error handling here
	return err
}
