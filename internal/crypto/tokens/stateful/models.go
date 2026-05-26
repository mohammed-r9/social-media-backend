package stateful

import (
	"time"

	"github.com/google/uuid"
)

type SessionTokens struct {
	RefreshToken string
	CsrfToken    string
}

type SessionTokenHashes struct {
	RefreshHash string
	CsrfHash    string
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

const (
	PASSWORD_RESET_TTL     = time.Minute * 30
	EMAIL_VERIFICATION_TTL = time.Hour * 2
	REFRESH_TTL            = time.Hour * 24 * 30
)

type ShortLivedToken struct {
	UserID    uuid.UUID  `json:"user_id"`
	Scope     TokenScope `json:"scope"`
	ExpiresAt time.Time  `json:"expires_at"`
	Ttl       time.Duration
	Raw       string
	Hash      string
}
