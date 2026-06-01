package rdrepo

import (
	"social-media-backend/internal/auth"
	"time"

	"github.com/google/uuid"
)

type StoredToken struct {
	UserID    uuid.UUID       `json:"user_id"`
	Scope     auth.TokenScope `json:"scope"`
	ExpiresAt time.Time       `json:"expires_at"`
}
