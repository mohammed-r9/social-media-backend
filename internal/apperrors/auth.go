package apperrors

import "net/http"

const (
	InvalidCredentials AuthError = iota
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

func (e AuthError) Error() string {
	switch e {
	case InvalidCredentials:
		return "invalid credentials"
	case UnverifiedUserEmail:
		return "unverified user email"
	case InvalidOldPassword:
		return "invalid user old password"
	case TokenMissmatch:
		return "token missmatch"
	case InvalidToken:
		return "invalid token"
	case EmptyRefreshToken:
		return "refresh token is empty"
	case EmptyCsrfToken:
		return "csrf token is empty"
	case TokenNotFound:
		return "token not found"
	case ExpiredToken:
		return "token expired"
	case MissingCsrf:
		return "missing csrf token"
	case MissingRefreshToken:
		return "missing refresh token"
	case MissingSessionID:
		return "missing session id"
	case MissingClaim:
		return "missing claim in jwt"
	case InvalidPassword:
		return "invalid user password"
	case MissingAuthorizationHeader:
		return "missing authorization header"
	default:
		return "unhandled auth error"
	}
}

func (e AuthError) Code() string {
	switch e {
	case InvalidCredentials:
		return "invalid_credentials"
	case UnverifiedUserEmail:
		return "email_not_verified"
	case InvalidOldPassword:
		return "invalid_old_password"
	case TokenMissmatch:
		return "token_mismatch"
	case InvalidToken:
		return "invalid_token"
	case EmptyRefreshToken:
		return "empty_refresh_token"
	case EmptyCsrfToken:
		return "empty_csrf_token"
	case TokenNotFound:
		return "token_not_found"
	case ExpiredToken:
		return "expired_token"
	case MissingCsrf:
		return "missing_csrf_token"
	case MissingRefreshToken:
		return "missing_refresh_token"
	case MissingSessionID:
		return "missing_session_id"
	case MissingClaim:
		return "missing_claim_in_jwt"
	case InvalidPassword:
		return "invalid_password"
	case MissingAuthorizationHeader:
		return "missing_authorization_header"
	default:
		return "unhandled_auth_error"
	}
}

func (e AuthError) Status() int {
	switch e {
	case InvalidCredentials:
		return http.StatusUnauthorized
	case UnverifiedUserEmail:
		return http.StatusForbidden
	case InvalidOldPassword:
		return http.StatusBadRequest
	case TokenMissmatch:
		return http.StatusUnauthorized
	case InvalidToken:
		return http.StatusUnauthorized
	case EmptyRefreshToken:
		return http.StatusBadRequest
	case EmptyCsrfToken:
		return http.StatusBadRequest
	case TokenNotFound:
		return http.StatusUnauthorized
	case ExpiredToken:
		return http.StatusUnauthorized
	case MissingCsrf:
		return http.StatusBadRequest
	case MissingRefreshToken:
		return http.StatusUnauthorized
	case MissingSessionID:
		return http.StatusUnauthorized
	case MissingClaim:
		return http.StatusUnauthorized
	case InvalidPassword:
		return http.StatusBadRequest
	case MissingAuthorizationHeader:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
