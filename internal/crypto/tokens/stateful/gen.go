package stateful

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

func NewShortLivedToken(userID uuid.UUID, scope TokenScope, ttl time.Duration) ShortLivedToken {
	raw := GenerateOpaqueToken(32)
	hash := hashToken(raw)

	return ShortLivedToken{
		UserID: userID,
		Scope:  scope,
		Ttl:    ttl,
		Raw:    raw,
		Hash:   hash,
	}
}
