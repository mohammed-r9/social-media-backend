package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

func GenerateOpaqueToken(length int) string {
	randomByte := make([]byte, length)
	rand.Read(randomByte)
	return hex.EncodeToString(randomByte)
}

func newShortLivedToken(userID uuid.UUID, scope TokenScope, ttl time.Duration) ShortLivedToken {
	raw := GenerateOpaqueToken(32)
	hash := HashToken(raw)

	return ShortLivedToken{
		UserID:    userID,
		Scope:     scope,
		Ttl:       ttl,
		Raw:       raw,
		Hash:      hash,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func GenerateEmailVerificationToken(userID uuid.UUID) ShortLivedToken {
	return newShortLivedToken(userID, ScopeEmailVerification, EMAIL_VERIFICATION_TTL)
}

func GeneratePasswordResetToken(userID uuid.UUID) ShortLivedToken {
	return newShortLivedToken(userID, ScopePasswordReset, PASSWORD_RESET_TTL)
}

func GenerateSessionID() string {
	return GenerateOpaqueToken(24)
}

func GenerateSessionTokens() SessionTokens {
	return SessionTokens{
		RefreshToken: GenerateOpaqueToken(32),
		CsrfToken:    GenerateOpaqueToken(32),
	}
}
