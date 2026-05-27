package tokens

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMissingClaim = errors.New("missing claim in jwt")
)

const (
	PASSWORD_RESET_TTL     = time.Minute * 30
	EMAIL_VERIFICATION_TTL = time.Hour * 2
	REFRESH_TTL            = time.Hour * 24 * 30
	ACCESS_TOKEN_TTL       = time.Minute * 15
)

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
	IsEmailVerified bool
	// maybe I need more fields? idk
}

// HashTokens requires both RefreshToken and CsrfToken to be set before being called,
func (t *SessionTokens) ToHash() SessionTokenHashes {
	return SessionTokenHashes{
		RefreshHash: HashToken(t.RefreshToken),
		CsrfHash:    HashToken(t.CsrfToken),
	}
}

func GenerateSessionTokens() SessionTokens {
	return SessionTokens{
		RefreshToken: GenerateOpaqueToken(32),
		CsrfToken:    GenerateOpaqueToken(32),
	}
}

type TokenScope string

const (
	ScopePasswordReset     TokenScope = "password_reset"
	ScopeEmailVerification TokenScope = "email_verification"
)

type ShortLivedToken struct {
	UserID    uuid.UUID  `json:"user_id"`
	Scope     TokenScope `json:"scope"`
	ExpiresAt time.Time  `json:"expires_at"`
	Ttl       time.Duration
	Raw       string
	Hash      string
}
