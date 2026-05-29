package network

import "errors"

var (
	errMissingAuthModeHeader      = errors.New("missing X-Auth-Mode header")
	errInvalidAuthModeHeader      = errors.New("invalid X-Auth-Mode header")
	errMissingAccessToken         = errors.New("missing access token")
	errInvalidAccessToken         = errors.New("invalid access token")
	errMissingAuthorizationHeader = errors.New("missing authorization header")
)
