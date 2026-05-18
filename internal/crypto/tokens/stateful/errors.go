package stateful

import "errors"

var (
	ErrEmptyRefreshToken = errors.New("refresh token is empty")
	ErrEmptyCsrfToken    = errors.New("csrf token is empty")
)
