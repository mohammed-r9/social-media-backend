package domain

import (
	"social-media-backend/internal/apperrors"
	"time"

	"github.com/google/uuid"
)

// TODO: I need to find a way to embed some user field here without breaking the design
type Session struct {
	ID               string
	UserID           uuid.UUID
	DeviceName       string
	RefreshTokenHash string
	CsrfTokenHash    string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
}

type CreateSessionParams struct {
	ID               string
	UserID           uuid.UUID
	DeviceName       string
	RefreshTokenHash string
	CsrfTokenHash    string
	ExpiresAt        time.Time
}

type GetUserSessionParams struct {
	ID     string
	UserID uuid.UUID
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	CsrfToken    string
	SessionID    string
}

func (s *Session) ValidateSession() error {
	if time.Now().After(s.ExpiresAt) {
		return apperrors.SessionExpired
	}

	if s.RevokedAt != nil {
		return apperrors.SessionRevoked
	}

	return nil
}
