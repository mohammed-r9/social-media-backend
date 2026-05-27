package network

import "errors"

var (
	errMissingAuthModeHeader = errors.New("missing X-Auth-Mode header")
	errInvalidAuthModeHeader = errors.New("invalid X-Auth-Mode header")
)
