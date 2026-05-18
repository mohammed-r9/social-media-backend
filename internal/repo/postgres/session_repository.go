package postgres

import (
	"context"
	"social-media-backend/internal/adapters/sqlc/db"
	"social-media-backend/internal/domain"

	"github.com/google/uuid"
)

func (r *PostgresRepo) CreateSession(ctx context.Context, params domain.CreateSessionParams) (uuid.UUID, error) {
	sessionID, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		ID:               params.ID,
		UserID:           params.UserID,
		RefreshTokenHash: params.RefreshTokenHash,
		CsrfTokenHash:    params.CsrfTokenHash,
		DeviceName:       params.DeviceName,
		ExpiresAt:        params.ExpiresAt,
	})

	// TODO: add error handling
	_ = err

	return sessionID, nil
}
