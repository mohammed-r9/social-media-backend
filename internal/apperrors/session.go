package apperrors

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
