package domain

import (
	"errors"
	"social-media-backend/internal/crypto/tokens/stateful"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type StoreTokenParam struct {
	Token stateful.ShortLivedToken
}
