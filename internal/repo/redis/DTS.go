package rdrepo

import (
	"social-media-backend/internal/crypto/tokens/stateful"
	"time"

	"github.com/google/uuid"
)

type StoredToken struct {
	UserID    uuid.UUID           `json:"user_id"`
	Scope     stateful.TokenScope `json:"scope"`
	ExpiresAt time.Time           `json:"expires_at"`
}
