package apperrors

import "net/http"

const (
	SessionExpired SessionError = iota
	SessionRevoked
	SessionAlreadyExists
	SessionNotFound
)

func (e SessionError) Error() string {
	switch e {
	case SessionExpired:
		return "session expired"
	case SessionRevoked:
		return "session revoked"
	case SessionAlreadyExists:
		return "session already exists"
	case SessionNotFound:
		return "session not found"
	default:
		return "unhandled session error"
	}
}

func (e SessionError) Code() string {
	switch e {
	case SessionExpired:
		return "session_expired"
	case SessionRevoked:
		return "session_revoked"
	case SessionAlreadyExists:
		return "session_exists"
	case SessionNotFound:
		return "session_not_found"
	default:
		return "unhandled_session_error"
	}
}

func (e SessionError) Status() int {
	switch e {
	case SessionExpired:
		return http.StatusUnauthorized
	case SessionRevoked:
		return http.StatusUnauthorized
	case SessionAlreadyExists:
		return http.StatusConflict
	case SessionNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
