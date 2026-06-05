package repo

import (
	"context"
	"database/sql"
	"errors"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type SessionsRepository interface {
	CreateSession(context.Context, domain.CreateSessionParams) (domain.Session, error)
	GetUserSession(context.Context, string) (UserSessionDTO, error)
}

var _ SessionsRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) CreateSession(ctx context.Context, params domain.CreateSessionParams) (domain.Session, error) {
	session, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		ID:               params.ID,
		UserID:           params.UserID,
		RefreshTokenHash: params.RefreshTokenHash,
		CsrfTokenHash:    params.CsrfTokenHash,
		DeviceName:       params.DeviceName,
		ExpiresAt:        params.ExpiresAt,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			// unique_violation
			case "23505":
				return domain.Session{}, domain.ErrSessionAlreadyExists

			// foreign_key_violation
			case "23503":
				return domain.Session{}, domain.ErrUserNotFound
			}
		}

		return domain.Session{}, err
	}

	return domain.Session{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		CsrfTokenHash:    session.CsrfTokenHash,
		DeviceName:       session.DeviceName,
		ExpiresAt:        session.ExpiresAt,
		RevokedAt:        session.RevokedAt,
	}, nil
}

func (r *PostgresRepo) GetUserSession(ctx context.Context, sessionID string) (UserSessionDTO, error) {
	session, err := r.q.GetUserSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserSessionDTO{}, domain.ErrSessionNotFound
		}

		return UserSessionDTO{}, err
	}

	return UserSessionDTO{
		Session: domain.Session{
			ID:               session.ID,
			UserID:           session.UserID,
			RefreshTokenHash: session.RefreshTokenHash,
			CsrfTokenHash:    session.CsrfTokenHash,
			DeviceName:       session.DeviceName,
			ExpiresAt:        session.ExpiresAt,
			RevokedAt:        session.RevokedAt,
		},
		User: struct{ VerifiedAt *time.Time }{
			VerifiedAt: session.VerifiedAt,
		},
	}, nil
}
