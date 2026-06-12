package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	UserID          uuid.UUID `json:"-"`
	Username        string    `json:"username"`
	IsEmailVerified bool      `json:"is_email_verified"`
	jwt.RegisteredClaims
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
