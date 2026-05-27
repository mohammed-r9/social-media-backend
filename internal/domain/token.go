package domain

import (
	"errors"
	"social-media-backend/internal/crypto/tokens/stateful"
)

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrExpiredToken        = errors.New("token not found")
	ErrMissingCsrf         = errors.New("missing csrf token")
	ErrMissingRefreshToken = errors.New("missing refresh token")
	ErrMissingSessionID    = errors.New("missing session id")
	ErrTokenMismatch       = errors.New("token mismatch")
)

type StoreTokenParam struct {
	Token stateful.ShortLivedToken
}
