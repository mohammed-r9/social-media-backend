package repo

import (
	"social-media-backend/internal/auth"
	"social-media-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

type UserSessionDTO struct {
	Session domain.Session
	User    struct {
		VerifiedAt *time.Time
	}
}

type StoredToken struct {
	UserID    uuid.UUID       `json:"user_id"`
	Scope     auth.TokenScope `json:"scope"`
	ExpiresAt time.Time       `json:"expires_at"`
}
