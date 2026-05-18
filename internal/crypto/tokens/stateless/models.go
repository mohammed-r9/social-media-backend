package stateless

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMissingClaim = errors.New("missing claim in jwt")
)

const (
	ACCESS_TOKEN_TTL = time.Minute * 15
)

type AccessTokenClaims struct {
	UserID          uuid.UUID
	IsEmailVerified bool
	// maybe I need more fields? idk
}

