package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionExpired = errors.New("session is expired")
	ErrSessionRevoked = errors.New("session is revoked")
)

type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	DeviceName       string
	RefreshTokenHash string
	CsrfTokenHash    string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
}

type CreateSessionParams struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	DeviceName       string
	RefreshTokenHash string
	CsrfTokenHash    string
	ExpiresAt        time.Time
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	CsrfToken    string
}

type SessionsRepository interface {
	CreateSession(context.Context, CreateSessionParams) (Session, error)
}

func (s *Session) ValidateSession() error {
	if time.Now().After(s.ExpiresAt) {
		return ErrSessionExpired
	}

	if s.RevokedAt != nil {
		return ErrSessionRevoked
	}

	return nil
}
