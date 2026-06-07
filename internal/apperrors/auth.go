package apperrors

import "net/http"

const (
	InvalidCredentials authError = iota
	UnverifiedUserEmail
	InvalidOldPassword
	TokenMissmatch
	InvalidToken
	EmptyRefreshToken
	EmptyCsrfToken
	TokenNotFound
	ExpiredToken
	MissingCsrf
	MissingRefreshToken
	MissingSessionID
	MissingClaim
	InvalidPassword
	MissingToken
	MissingAuthorizationHeader
)

var authErrorTable = [...]errorInfo{
	InvalidCredentials: {
		message: "invalid credentials",
		code:    "invalid_credentials",
		status:  http.StatusUnauthorized,
	},
	UnverifiedUserEmail: {
		message: "unverified user email",
		code:    "email_not_verified",
		status:  http.StatusForbidden,
	},
	InvalidOldPassword: {
		message: "invalid user old password",
		code:    "invalid_old_password",
		status:  http.StatusBadRequest,
	},
	TokenMissmatch: {
		message: "token missmatch",
		code:    "token_mismatch",
		status:  http.StatusUnauthorized,
	},
	InvalidToken: {
		message: "invalid token",
		code:    "invalid_token",
		status:  http.StatusUnauthorized,
	},
	EmptyRefreshToken: {
		message: "refresh token is empty",
		code:    "empty_refresh_token",
		status:  http.StatusBadRequest,
	},
	EmptyCsrfToken: {
		message: "csrf token is empty",
		code:    "empty_csrf_token",
		status:  http.StatusBadRequest,
	},
	TokenNotFound: {
		message: "token not found",
		code:    "token_not_found",
		status:  http.StatusUnauthorized,
	},
	ExpiredToken: {
		message: "token expired",
		code:    "expired_token",
		status:  http.StatusUnauthorized,
	},
	MissingCsrf: {
		message: "missing csrf token",
		code:    "missing_csrf_token",
		status:  http.StatusBadRequest,
	},
	MissingRefreshToken: {
		message: "missing refresh token",
		code:    "missing_refresh_token",
		status:  http.StatusUnauthorized,
	},
	MissingSessionID: {
		message: "missing session id",
		code:    "missing_session_id",
		status:  http.StatusUnauthorized,
	},
	MissingClaim: {
		message: "missing claim in jwt",
		code:    "missing_claim_in_jwt",
		status:  http.StatusUnauthorized,
	},
	InvalidPassword: {
		message: "invalid user password",
		code:    "invalid_password",
		status:  http.StatusBadRequest,
	},
	MissingAuthorizationHeader: {
		message: "missing authorization header",
		code:    "missing_authorization_header",
		status:  http.StatusUnauthorized,
	},
}

func (e authError) Error() string {
	return e.info().message
}

func (e authError) Code() string {
	return e.info().code
}

func (e authError) Status() int {
	return e.info().status
}

func (e authError) info() errorInfo {
	if int(e) < 0 || int(e) >= len(authErrorTable) {
		return errorInfo{
			message: "unhandled auth error",
			code:    "unhandled_auth_error",
			status:  http.StatusInternalServerError,
		}
	}
	return authErrorTable[e]
}
