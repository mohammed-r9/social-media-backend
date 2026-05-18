package domain

import (
	"time"

	"github.com/google/uuid"
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

type SessionsRepository interface {
}
