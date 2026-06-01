package auth

import (
	"time"

	"github.com/google/uuid"
)

type TokenScope string

type SessionTokens struct {
	RefreshToken string
	CsrfToken    string
}

type SessionTokenHashes struct {
	RefreshHash string
	CsrfHash    string
}

type AccessTokenClaims struct {
	UserID          uuid.UUID
	Username        string
	IsEmailVerified bool
}

type ShortLivedToken struct {
	UserID    uuid.UUID  `json:"user_id"`
	Scope     TokenScope `json:"scope"`
	ExpiresAt time.Time  `json:"expires_at"`
	Ttl       time.Duration
	Raw       string
	Hash      string
}

type StoreTokenParam struct {
	Token ShortLivedToken
}

// HashTokens requires both RefreshToken and CsrfToken to be set before being called,
func (t *SessionTokens) ToHash() SessionTokenHashes {
	return SessionTokenHashes{
		RefreshHash: HashToken(t.RefreshToken),
		CsrfHash:    HashToken(t.CsrfToken),
	}
}
