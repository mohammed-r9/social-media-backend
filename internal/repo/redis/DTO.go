package rdrepo

import (
	"social-media-backend/internal/crypto/tokens"
	"time"

	"github.com/google/uuid"
)

type StoredToken struct {
	UserID    uuid.UUID         `json:"user_id"`
	Scope     tokens.TokenScope `json:"scope"`
	ExpiresAt time.Time         `json:"expires_at"`
}
