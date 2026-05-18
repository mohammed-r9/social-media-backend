package domain

import (
	"context"
	"errors"
	"social-media-backend/internal/crypto/tokens/stateful"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionExpired = errors.New("session is expired")
	ErrSessionRevoked = errors.New("session is revoked")
)

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

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	CsrfToken    string
	SessionID    string
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

func GenerateSessionID() string {
	return stateful.GenerateOpaqueToken(24)
}
