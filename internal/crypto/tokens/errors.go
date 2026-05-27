package tokens

import "errors"

var (
	ErrEmptyRefreshToken   = errors.New("refresh token is empty")
	ErrEmptyCsrfToken      = errors.New("csrf token is empty")
	ErrTokenNotFound       = errors.New("token not found")
	ErrExpiredToken        = errors.New("token not found")
	ErrMissingCsrf         = errors.New("missing csrf token")
	ErrMissingRefreshToken = errors.New("missing refresh token")
	ErrMissingSessionID    = errors.New("missing session id")
	ErrTokenMismatch       = errors.New("token mismatch")
	ErrInvalidToken        = errors.New("invalid token")
)
