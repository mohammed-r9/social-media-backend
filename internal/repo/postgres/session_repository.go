package postgres

import (
	"context"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"
)

var _ domain.SessionsRepository = (*PostgresRepo)(nil)

func (r *PostgresRepo) CreateSession(ctx context.Context, params domain.CreateSessionParams) (domain.Session, error) {
	session, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		ID:               params.ID,
		UserID:           params.UserID,
		RefreshTokenHash: params.RefreshTokenHash,
		CsrfTokenHash:    params.CsrfTokenHash,
		DeviceName:       params.DeviceName,
		ExpiresAt:        params.ExpiresAt,
	})

	// TODO: add error handling
	_ = err

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

func (r *PostgresRepo) GetUserSession(ctx context.Context, params domain.GetUserSessionParams) (domain.Session, error) {
	session, err := r.q.GetUserSession(ctx, db.GetUserSessionParams{
		ID:     params.ID,
		UserID: params.UserID,
	})

	// TODO: add error handling
	_ = err

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
