package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

// hashes the token and encodes the hash as a string
func HashToken(plainText string) string {
	hash := sha256.Sum256([]byte(plainText))
	return hex.EncodeToString(hash[:])
}

type CompareTokenToHashParams struct {
	PlainText  string
	StoredHash string
}

func CompareTokenToHash(params CompareTokenToHashParams) bool {
	tokenHash := HashToken(params.PlainText)
	return subtle.ConstantTimeCompare([]byte(tokenHash), []byte(params.StoredHash)) == 1
}

func RedisTokenKeyBuilder(tokenHash string, scope TokenScope) string {
	switch scope {
	case SCOPE_PASSWORD_RESET:
		return "prt:" + tokenHash
	case SCOPE_EMAIL_VERIFICATION:
		return "evt:" + tokenHash
	default:
		return ""
	}
}
